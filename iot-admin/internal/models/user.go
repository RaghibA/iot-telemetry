package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   string `gorm:"primaryKey"`
	Username string
	Password []byte
	Email    string
	APIKey   string
}
