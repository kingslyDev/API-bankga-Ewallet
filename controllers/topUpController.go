package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/models"
	"github.com/kingslyDev/API-bankga-Ewallet/utils"
	"gorm.io/gorm"
)

type TopUpRequest struct {
	Amount            float64 `json:"amount" binding:"required,min=10000"`
	Pin               string  `json:"pin" binding:"required,len=6"`
	PaymentMethodCode string  `json:"payment_method_code" binding:"required,oneof=gopay bca_va"`
}

type TopUpController struct {
	DB *gorm.DB
}

func (ctrl *TopUpController) TopUp(c *gin.Context) {
	// Ambil klaim dari JWT
	claims, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no claims found"})
		return
	}

	email, ok := claims.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid claims format"})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := ctrl.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user not found"})
		return
	}

	// Validasi input
	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari wallet berdasarkan user ID
	var wallet models.Wallet
	if err := ctrl.DB.Where("user_id = ?", user.ID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found"})
		return
	}

	// Validasi PIN
	if wallet.Pin == "" || wallet.Pin != req.Pin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
		return
	}

	// Cari TransactionType dan PaymentMethod
	var transactionType models.TransactionType
	if err := ctrl.DB.Where("code = ?", "top_up").First(&transactionType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction type not found"})
		return
	}

	var paymentMethod models.PaymentMethods
	if err := ctrl.DB.Where("code = ?", req.PaymentMethodCode).First(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment method not found"})
		return
	}

	// Mulai transaksi
	tx := ctrl.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate kode transaksi
	transactionCode := strings.ToUpper(utils.RandomString(10))

	// Buat objek transaksi
	desc := "Topup via " + paymentMethod.Name
	transaction := models.Transaction{
		UserID:            user.ID,
		PaymentMethodID:   paymentMethod.ID,
		TransactionTypeID: transactionType.ID,
		Amount:            req.Amount,
		TransactionCode:   transactionCode,
		Description:       &desc,
		Status:            "pending",
	}

	// Simpan transaksi ke database
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Panggil Midtrans API
	midtransParams := utils.BuildMidtransParams(transaction.TransactionCode, transaction.Amount, user)
	midtransResp, err := utils.CallMidtrans(midtransParams)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	// Commit transaksi
	tx.Commit()

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
		"midtrans":    midtransResp,
	})
}
