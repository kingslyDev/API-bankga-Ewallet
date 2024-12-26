package models

import "gorm.io/gorm"

type DataPlans struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Price float64 `json:"price" gorm:"type:decimal(10,2); not null"`
	OperatorCardID uint `json:"operator_card_id"`
	OperatorCard OperatorCard `json:"operator_card" gorm:"constraint:OnDelete:CASCADE;"`
}