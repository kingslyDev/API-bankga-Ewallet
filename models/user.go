package models

import (
	"time"
)

	type User struct {
		ID             uint      `gorm:"primaryKey"`
		Name           string    `gorm:"size:255;not null"`
		Email          string    `gorm:"size:255;unique;not null"`
		Username       string    `gorm:"size:255;unique;not null"`
		Password       string    `gorm:"not null"`
		Verified       *bool     `gorm:"default:false"`
		ProfilePicture string    `gorm:"size:255"`
		KTP            string    `gorm:"size:255"`
		RememberToken  string    `gorm:"size:255"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
