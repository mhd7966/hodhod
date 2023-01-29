package models

import (

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID          int          `json:"id" `
	UserID      int          `json:"user_id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	PhoneNumber string       `json:"phone_number" binding:"required"`
	Email       string       `json:"email" binding:"required"`
}
