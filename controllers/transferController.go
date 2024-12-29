package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/models"
	"gorm.io/gorm"
)

type TransferController struct {
	DB *gorm.DB
}

// Input struct for transfer request
type TransferRequest struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Pin     string  `json:"pin" binding:"required,len=6"`
	SendTo  string  `json:"send_to" binding:"required"`
}

// Transfer handles fund transfer between users
func (ctrl *TransferController) Transfer(c *gin.Context) {
	var input TransferRequest

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Get sender information (authenticated user)
	sender, ok := c.MustGet("user").(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Find receiver by username or card number
	receiver, err := ctrl.findReceiver(input.SendTo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Ensure sender and receiver are not the same
	if sender.ID == receiver.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot transfer to yourself"})
		return
	}

	// Verify sender's PIN
	if !ctrl.verifyPin(sender, input.Pin) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PIN"})
		return
	}

	// Start database transaction
	tx := ctrl.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	// Process the transfer
	transactionCode, err := ctrl.processTransfer(tx, sender, receiver, input.Amount)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete the transaction"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful", "transaction_code": transactionCode})
}

// Find receiver by username or card number
func (ctrl *TransferController) findReceiver(identifier string) (*models.User, error) {
	var receiver models.User
	err := ctrl.DB.Joins("JOIN wallets ON wallets.user_id = users.id").
		Where("users.username = ? OR wallets.card_number = ?", identifier, identifier).
		First(&receiver).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receiver not found")
		}
		return nil, err
	}
	return &receiver, nil
}

// Verify sender's PIN
func (ctrl *TransferController) verifyPin(user *models.User, pin string) bool {
	var wallet models.Wallet
	if err := ctrl.DB.Where("user_id = ?", user.ID).First(&wallet).Error; err != nil {
		log.Println("Wallet not found for user:", user.ID)
		return false
	}

	// Verifikasi PIN pada wallet pengirim
	return wallet.Pin == pin
}

// Process the transfer
func (ctrl *TransferController) processTransfer(tx *gorm.DB, sender, receiver *models.User, amount float64) (string, error) {
	var senderWallet, receiverWallet models.Wallet

	// Get sender wallet
	if err := tx.Where("user_id = ?", sender.ID).First(&senderWallet).Error; err != nil {
		return "", errors.New("sender wallet not found")
	}

	// Check sender's balance
	if senderWallet.Balance < amount {
		return "", errors.New("insufficient balance")
	}

	// Get receiver wallet
	if err := tx.Where("user_id = ?", receiver.ID).First(&receiverWallet).Error; err != nil {
		return "", errors.New("receiver wallet not found")
	}

	// Generate unique transaction code
	transactionCode := fmt.Sprintf("TRF-%s-%d", time.Now().Format("20060102150405"), sender.ID)

	// Deduct amount from sender
	senderWallet.Balance -= amount
	if err := tx.Save(&senderWallet).Error; err != nil {
		return "", errors.New("failed to update sender balance")
	}

	// Add amount to receiver
	receiverWallet.Balance += amount
	if err := tx.Save(&receiverWallet).Error; err != nil {
		return "", errors.New("failed to update receiver balance")
	}

	// Create transaction records
	if err := ctrl.createTransactions(tx, sender, receiver, amount, transactionCode); err != nil {
		return "", err
	}

	// Create transaction history
	history := models.TransactionHistory{
		SenderID:       sender.ID,
		ReceiverID:     receiver.ID,
		TransactionCode: transactionCode,
	}
	if err := tx.Create(&history).Error; err != nil {
		return "", errors.New("failed to create transaction history")
	}

	return transactionCode, nil
}

// Create transaction records for both sender and receiver
func (ctrl *TransferController) createTransactions(tx *gorm.DB, sender, receiver *models.User, amount float64, transactionCode string) error {
	var transferType, receiveType models.TransactionType
	var paymentMethod models.PaymentMethods

	// Get transaction types
	if err := tx.Where("code = ? AND action = ?", "transfer", "cr").First(&transferType).Error; err != nil {
		return errors.New("transfer transaction type not found")
	}
	if err := tx.Where("code = ? AND action = ?", "receive", "dr").First(&receiveType).Error; err != nil {
		return errors.New("receive transaction type not found")
	}

	// Get payment method
	if err := tx.Where("code = ?", "bwa").First(&paymentMethod).Error; err != nil {
		return errors.New("payment method not found")
	}

	// Create sender transaction
	senderTransaction := models.Transaction{
		UserID:           sender.ID,
		TransactionTypeID: transferType.ID,
		Description:      ptr(fmt.Sprintf("Transfer to %s", receiver.Username)),
		Amount:           amount,
		TransactionCode:  transactionCode,
		Status:           "success",
		PaymentMethodID:  paymentMethod.ID,
	}
	if err := tx.Create(&senderTransaction).Error; err != nil {
		return errors.New("failed to create sender transaction")
	}

	// Create receiver transaction
	receiverTransaction := models.Transaction{
		UserID:           receiver.ID,
		TransactionTypeID: receiveType.ID,
		Description:      ptr(fmt.Sprintf("Received from %s", sender.Username)),
		Amount:           amount,
		TransactionCode:  transactionCode,
		Status:           "success",
		PaymentMethodID:  paymentMethod.ID,
	}
	if err := tx.Create(&receiverTransaction).Error; err != nil {
		return errors.New("failed to create receiver transaction")
	}

	return nil
}

// Helper function to return pointer of string
func ptr(s string) *string {
	return &s
}
