package models

import (
	"gorm.io/gorm"
)

type TransactionType struct {
    gorm.Model
    Name      string         `json:"name" gorm:"size:255;not null"`
    Code      string         `json:"code" gorm:"size:255;unique;not null"`
    Action    string         `json:"action" gorm:"type:action_type;not null"`  // Menggunakan tipe enum yang sudah ada
    Thumbnail string         `json:"thumbnail" gorm:"size:255;not null"`
    DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
