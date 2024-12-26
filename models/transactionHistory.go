package models

import "gorm.io/gorm"

type TransactionHistory struct {
	gorm.Model
	SenderID uint `json:"sender_id"`
	Sender User `json:"sender" gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;"`
	ReceiverID uint `json:"receiver_id"`
	Receiver User `json:"receiver" gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE;"`
	TransactionCode string `json:"transaction_code" gorm:"size:255;not null"`
}