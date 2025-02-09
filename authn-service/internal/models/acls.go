package models

import "gorm.io/gorm"

type KafkaACL struct {
	gorm.Model
	UserID    string
	DeviceID  string
	APIKey    string
	TopicName string
	Read      bool
	Write     bool
}
