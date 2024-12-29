package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/kingslyDev/API-bankga-Ewallet/controllers" // sesuaikan nama import
	"github.com/kingslyDev/API-bankga-Ewallet/controllers/middleware"
	"gorm.io/gorm"
)


func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
  
    authController := &controller.AuthController{DB: db}
	topUpController := &controller.TopUpController{DB: db}
	webhookController := &controller.WebhookController{DB: db}
	log.Println("Registering routes...")
	router.POST("/api/register", authController.Register)
	router.POST("/api/login", authController.Login)
	
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", authController.GetProfile)
		protected.POST("/topup", topUpController.TopUp)
	}
	router.POST("/api/webhook/midtrans", webhookController.UpdateTransaction)
}
