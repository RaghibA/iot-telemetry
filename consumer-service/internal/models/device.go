package models

import "gorm.io/gorm"

/**
 * Device Schema
 *
 * Init on device creation
 */
type Device struct {
	gorm.Model
	DeviceName string
	DeviceID   string
	UserID     string
	TopicName  string
}
