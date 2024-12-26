package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/controllers"
	"gorm.io/gorm"
)

// RegisterRoutes akan mendefinisikan semua route aplikasi
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Membuat instance authController
	authController := &controllers.AuthController{DB: db}

	// Endpoint untuk registrasi
	router.POST("/api/register", authController.Register)
}
