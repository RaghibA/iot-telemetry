package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	DeviceName string
	DeviceID   string
	UserID     string
	TopicName  string
}
