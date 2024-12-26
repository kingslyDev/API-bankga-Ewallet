package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/kingslyDev/API-bankga-Ewallet/controllers" // sesuaikan nama import
	"gorm.io/gorm"
)

// RegisterRoutes akan mendefinisikan semua route aplikasi
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    // Membuat instance authController
    authController := &controller.AuthController{DB: db}

    // Log untuk debugging
    log.Println("Registering routes...")

    // Endpoint untuk registrasi
    router.POST("/api/register", authController.Register)
    
    // Debugging rute
    log.Println("Routes registered:")
    for _, route := range router.Routes() {
        log.Printf("Route: %s %s", route.Method, route.Path)
    }
}
