package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == ""{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error" : "Authorization Header invalid"})
			ctx.Abort()
			return

		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		if tokenString == authHeader{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error" : "Invalid Token Format"})
			ctx.Abort()
			return
		}

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized,gin.H{"error" : "Invalid or expires token"})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Next()

	}
}