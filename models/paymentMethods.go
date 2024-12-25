package models

import "gorm.io/gorm"

type PaymentStatus string
const (
	Active PaymentStatus = "active"
	Inactive PaymentStatus = "inactive"
)


type PaymentMethods struct {
	gorm.Model
	Name string `json:"name" gorm:"size:255; not null"`
	Code string `json:"code" gorm:"size:255; not null"`
	Thumbnail *string `json:"thumbnail" gorm:"size:255"`
	Status PaymentStatus `json:"status" gorm:"type:payment_status; not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}