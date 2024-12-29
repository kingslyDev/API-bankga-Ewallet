package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/kingslyDev/API-bankga-Ewallet/controllers" // sesuaikan nama import
	"github.com/kingslyDev/API-bankga-Ewallet/controllers/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Membuat instance controller dengan db
	authController := &controller.AuthController{DB: db}
	topUpController := &controller.TopUpController{DB: db}
	webhookController := &controller.WebhookController{DB: db}
	transferController := &controller.TransferController{DB: db}

	log.Println("Registering routes...")

	// Route untuk register dan login
	router.POST("/api/register", authController.Register)
	router.POST("/api/login", authController.Login)

	// Group untuk route yang membutuhkan autentikasi
	protected := router.Group("/api")
	// Menambahkan middleware JWTAuthMiddleware yang sekarang menerima parameter db
	protected.Use(middleware.JWTAuthMiddleware(db))
	{
		// Protected routes
		protected.GET("/profile", authController.GetProfile)
		protected.POST("/topup", topUpController.TopUp)
		protected.POST("/transfer", transferController.Transfer)
	}

	// Route untuk webhook
	router.POST("/api/webhook/midtrans", webhookController.UpdateTransaction)
}
