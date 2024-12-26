package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kingslyDev/API-bankga-Ewallet/config"
	"github.com/kingslyDev/API-bankga-Ewallet/routers"
)

func main() {
	// Inisialisasi koneksi ke database
	config.ConnectDatabase() // Pastikan ini dipanggil untuk koneksi ke database

	// Periksa apakah koneksi berhasil
	if config.DB == nil {
		log.Fatal("Failed to connect to database")
		return
	}

	// Membuat instance Gin router
	r := gin.Default()

	// Daftarkan semua routes
	routers.RegisterRoutes(r, config.DB)

	// Menjalankan server pada port 8080
	if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server on port 8080: %v", err)
    } else {
        log.Println("Server started on http://localhost:8080")
    }
    }
