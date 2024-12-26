package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/kingslyDev/API-bankga-Ewallet/controllers" // sesuaikan nama import
	"gorm.io/gorm"
)


func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
  
    authController := &controller.AuthController{DB: db}
    log.Println("Registering routes...")
    router.POST("/api/register", authController.Register)
    
    // Debugging rute
    log.Println("Routes registered:")
    for _, route := range router.Routes() {
        log.Printf("Route: %s %s", route.Method, route.Path)
    }
}
