package controller

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/kingslyDev/API-bankga-Ewallet/config"
	"github.com/kingslyDev/API-bankga-Ewallet/models"
	"github.com/kingslyDev/API-bankga-Ewallet/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	User        *models.User `json:"user"`
	Token       string       `json:"token"`
	TokenType   string       `json:"token_type"`
	TokenExpires time.Time   `json:"token_expires"`
}

type AuthController struct {
	DB *gorm.DB
}

// GetProfile retrieves the profile of the currently authenticated user
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Login authenticates a user and returns a JWT token
func (ctrl *AuthController) Login(c *gin.Context) {
	var loginReq LoginRequest

	// Validate input
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		}
		return
	}

	// Check password hash
	if !utils.CheckPasswordHash(loginReq.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Generate JWT token
	token, expires, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	response := LoginResponse{
		User:         &user,
		Token:        token,
		TokenType:    "Bearer",
		TokenExpires: expires,
	}

	c.JSON(http.StatusOK, response)
}

// Register registers a new user and creates a wallet
func (ctrl *AuthController) Register(c *gin.Context) {
	var data struct {
		Name           string `json:"name" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		Username       string `json:"username" binding:"required,alphanum"`
		Password       string `json:"password" binding:"required,min=8"`
		Pin            string `json:"pin" binding:"required,len=6"`
		ProfilePicture string `json:"profile_picture"`
		KTP            string `json:"ktp"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate PIN length
	if len(data.Pin) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be 6 digits"})
		return
	}

	// Check if email already exists
	if exists := isEmailExists(data.Email); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	// Check if username already exists
	if exists := isUsernameExists(data.Username); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already in use"})
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Upload profile picture and KTP
	profilePicturePath := utils.UploadBase64Image(data.ProfilePicture, "profile_pictures")
	ktpPath := utils.UploadBase64Image(data.KTP, "ktp")

	// Create user
	user := models.User{
		Name:           data.Name,
		Email:          data.Email,
		Username:       data.Username,
		Password:       utils.HashPassword(data.Password),
		ProfilePicture: profilePicturePath,
		KTP:            ktpPath,
		Verified:       boolPtr(ktpPath != ""),
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create wallet
	wallet := models.Wallet{
		UserID:     user.ID,
		CardNumber: generateCardNumber(),
		Pin:        data.Pin,
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Commit transaction
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// Helper function to check if email exists
func isEmailExists(email string) bool {
	var user models.User
	return config.DB.Where("email = ?", email).First(&user).Error == nil
}

// Helper function to check if username exists
func isUsernameExists(username string) bool {
	var user models.User
	return config.DB.Where("username = ?", username).First(&user).Error == nil
}

// Generate a random 16-digit card number
func generateCardNumber() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%016x", b)
}

// Return a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}