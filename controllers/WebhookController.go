package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/models"
	"github.com/kingslyDev/API-bankga-Ewallet/utils"
	"gorm.io/gorm"
)

type WebhookController struct {
	DB *gorm.DB
}

func (ctrl *WebhookController) UpdateTransaction(c *gin.Context) {
	// Cek Midtrans server key
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if midtransServerKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Midtrans server key is not set"})
		return
	}

	// Parsing notifikasi Midtrans
	notif, err := utils.ParseMidtransNotification(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Midtrans Notification"})
		return
	}

	transactionStatus := notif.TransactionStatus
	transactionCode := notif.OrderID
	fraudStatus := notif.FraudStatus

	// Mulai transaksi database
	tx := ctrl.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	}()

	var status string

	// Tentukan status berdasarkan notifikasi
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			status = "success"
		}
	case "settlement":
		status = "success"
	case "cancel", "deny", "expire":
		status = "failed"
	case "pending":
		status = "pending"
	default:
		status = "unknown"
	}

	// Cari transaksi berdasarkan transaction_code
	var transaction models.Transaction
	if err := tx.Where("transaction_code = ?", transactionCode).First(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Update status transaksi jika belum sukses
	if transaction.Status != "success" {
		transaction.Status = status
		if err := tx.Save(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
			return
		}

		// Jika status sukses, tambahkan saldo ke wallet pengguna
		if status == "success" {
			if err := tx.Model(&models.Wallet{}).
				Where("user_id = ?", transaction.UserID).
				Update("balance", gorm.Expr("balance + ?", transaction.Amount)).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
				return
			}
		}
	}

	// Commit transaksi database
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}
