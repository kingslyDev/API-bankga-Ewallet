package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/kingslyDev/API-bankga-Ewallet/controllers"
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

    // Route untuk register dan login (tanpa autentikasi)
    router.POST("/api/register", authController.Register)
    router.POST("/api/login", authController.Login)

    // Group untuk route yang membutuhkan autentikasi
    protected := router.Group("/api")
    protected.Use(middleware.JWTAuthMiddleware()) // Middleware tanpa db
    {
        // Protected routes
        protected.GET("/profile", authController.GetProfile)
        protected.POST("/topup", topUpController.TopUp)
        protected.POST("/transfer", transferController.Transfer)
    }

    // Route untuk webhook (tidak perlu autentikasi)
    router.POST("/api/webhook/midtrans", webhookController.UpdateTransaction)
}
