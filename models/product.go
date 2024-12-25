package models

import "gorm.io/gorm"

type Product struct{
	gorm.Model
	Name string `json:"name" gorm:"size:255; not null;"`
	Thumbnail *string `json:"thumbnail" gorm:"size:255; not null"`
	Price float64 `json:"price" gorm:"type:decimal(10,2); not null"`
	Status string `json:"status" gorm:"type:product_status; not null"`
	Description *string `json:"description" gorm:"type:text"`
}