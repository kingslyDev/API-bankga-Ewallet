package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
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
		log.Println("Received token:", tokenString) // Log untuk debugging

		// Verifikasi token menggunakan ParseJWT
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			log.Println("JWT Error:", err) // Log kesalahan untuk debugging
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Set klaim email di context untuk digunakan di handler
		ctx.Set("email", claims.Email)
		log.Println("Email set in context:", claims.Email) // Log untuk memastikan email berhasil disimpan

		// Lanjutkan ke handler berikutnya
		ctx.Next()
	}
}
