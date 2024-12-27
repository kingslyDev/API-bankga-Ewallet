package middleware

import (
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

		// Menghapus "Bearer " dan spasi setelahnya
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Format"})
			ctx.Abort()
			return
		}

		// Verifikasi token menggunakan ParseJWT
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Set email di context untuk digunakan di handler
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}
