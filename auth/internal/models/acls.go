package models

import "gorm.io/gorm"

/**
 * KafkaACL Schema
 *
 * Init at device creation, but required update in
 * auth service if new api key is generated
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
