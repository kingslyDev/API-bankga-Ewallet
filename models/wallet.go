package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Balance float64 `json:"balance" gorm:"type:decimal(10,2); not null"`
	Pin string `json:"pin"`
	UserID uint `json:"user_id"`
	User User `json:"user" gorm:"constraint:OnDelete:CASCADE;"`
	CardNumber string `json:"card_number" gorm:"unique; not null"`
}