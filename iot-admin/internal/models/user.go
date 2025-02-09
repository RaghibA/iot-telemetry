package models

import "gorm.io/gorm"

/**
 * User Schema
 *
 * used by admin service for user data
 * associated with api keys & devices
 */
type User struct {
	gorm.Model
	UserID   string `gorm:"primaryKey"`
	Username string
	Password []byte
	Email    string
	APIKey   string
}
