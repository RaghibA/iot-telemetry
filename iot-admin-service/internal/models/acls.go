package models

import "gorm.io/gorm"

/**
 * KafkaACL Schema
 *
 * Init at device creation
 */
type KafkaACL struct {
	gorm.Model
	UserID    string
	DeviceID  string
	APIKey    string
	TopicName string
	Read      bool
	Write     bool
}
