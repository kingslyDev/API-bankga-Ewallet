package main

import (
	"github.com/kingslyDev/API-bankga-Ewallet/config"
	"github.com/kingslyDev/API-bankga-Ewallet/models"
)

func main() {
	config.ConnectDatabase()
    config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.Wallet{})
    config.DB.AutoMigrate(&models.TransactionType{})
    config.DB.AutoMigrate(models.Product{})
}
