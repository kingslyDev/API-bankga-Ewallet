package models

import (
	"gorm.io/gorm"
)

type Transaction struct{
	gorm.Model
	UserID uint `json:"user_id"`
	User User `json:"user" gorm:"constraint:OnDelete:CASCADE;"`
	TransactionTypeID uint `json:"transaction_type_id"`
	TransactionType TransactionType `json:"transaction_type" gorm:"constraint:OnDelete:CASCADE;"`
	PaymentMethodID uint `json:"payment_method_id"`
	PaymentMethod PaymentMethods `json:"payment_method" gorm:"constraint:OnDelete:CASCADE;"`
	ProductID *uint `json:"product_id"`
	Product *Product `json:"product" gorm:"constraint:OnDelete:SET NULL;"`
	Amount float64 `json:"amount" gorm:"type:decimal(10,2);not null"`
	TransactionCode string `json:"transaction_code" gorm:"size:255;not null"` 
	Description *string `json:"description" gorm:"type:text"`
	Status string `json:"status" gorm:"size:255;not null"`
}