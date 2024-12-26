package models

import "gorm.io/gorm"

type AdminUser struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
}