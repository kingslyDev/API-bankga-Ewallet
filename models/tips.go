package models

import "gorm.io/gorm"

type Tips struct {
	gorm.Model
	Title string `json:"title"`
	Url string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}