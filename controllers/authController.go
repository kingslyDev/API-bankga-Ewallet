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

type LoginRequest struct 
{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`

}


type LoginResponse struct 
{
	User *models.User `json:"user"`
	Token string `json:"token"`
	TokenType string `json:"token_type"`
	TokenExpires time.Time `json:"token_expires"`
}

type AuthController struct {
	DB *gorm.DB
}


func (ctrl *AuthController) Login(c *gin.Context){
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error" : "Invalid Input"})
		return
	}

	// query email utk login
	var user models.User
	if err := config.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			c.JSON(http.StatusUnauthorized, gin.H{"error" : "Your Email / Password is wrong"})

		} else {
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "server down"})
		}
		return
	
	}

	// fungsi checkHashpassword
	if !utils.CheckPasswordHash(loginReq.Password, user.Password){
		c.JSON(http.StatusUnauthorized, gin.H{"error" : "Your Email / Password is incorrect"})
		return
	}

	token, expires, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Server Down"})
		return
	}

	response := LoginResponse {
		User: &user,
		Token: token,
		TokenType: "Bearer",
		TokenExpires: expires,
	}

	c.JSON(http.StatusOK,response)

}

// Register untuk registrasi user baru
func (ctrl *AuthController) Register(c *gin.Context) {
	var data struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		Pin             string `json:"pin"`
		ProfilePicture  string `json:"profile_picture"`
		KTP             string `json:"ktp"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validasi panjang PIN
	if len(data.Pin) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be 6 digits"})
		return
	}

	// Cek jika email sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ?", data.Email).First(&existingUser).Error; err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	// Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Upload gambar profile dan KTP
	profilePicturePath := utils.UploadBase64Image(data.ProfilePicture, "profile_pictures")
	ktpPath := utils.UploadBase64Image(data.KTP, "ktp")

	// Buat user baru
	user := models.User{
		Name:           data.Name,
		Email:          data.Email,
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

	// Buat wallet baru
	wallet := models.Wallet{
		UserID:     user.ID,
		CardNumber: generateCardNumber(),
		Pin:        data.Pin, // PIN diambil dari input
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Commit transaksi jika semua berhasil
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// Fungsi untuk menghasilkan nomor kartu secara acak
func generateCardNumber() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%016x", b)[:16]
}

// Fungsi untuk mengembalikan pointer bool
func boolPtr(b bool) *bool {
	return &b
}
