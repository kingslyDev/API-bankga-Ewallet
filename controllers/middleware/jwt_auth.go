package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/models"
	"github.com/kingslyDev/API-bankga-Ewallet/utils"
	"gorm.io/gorm"
)

func JWTAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Header missing"})
			ctx.Abort()
			return
		}

		// Validasi format token harus dimulai dengan "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Format"})
			ctx.Abort()
			return
		}

		// Menghapus "Bearer " dari header
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// Verifikasi token menggunakan ParseJWT
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			log.Println("JWT Error:", err) // Log kesalahan untuk debugging
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Ambil user berdasarkan email dari klaim token
		var user models.User
		if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
			log.Println("User not found:", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		// Set objek user ke dalam konteks
		ctx.Set("user", &user)

		// Lanjutkan ke handler berikutnya
		ctx.Next()
	}
}
