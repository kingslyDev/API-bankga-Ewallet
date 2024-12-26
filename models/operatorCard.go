package models

import (
	"gorm.io/gorm"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type OperatorCard struct {
	gorm.Model
	Name   string `json:"name" gorm:"size:255;not null"` 
	Status Status `json:"status" gorm:"type:status_type;not null;default:'active'"` 
	Thumbnail string `json:"thumbnail"`
}
