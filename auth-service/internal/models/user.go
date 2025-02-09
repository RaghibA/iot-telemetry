package models

import "gorm.io/gorm"

/**
 * User Schema
 *
 * Init on account creation
 */
type User struct {
	gorm.Model
	UserID   string `gorm:"primaryKey"`
	Username string
	Password []byte
	Email    string
	APIKey   string
}
