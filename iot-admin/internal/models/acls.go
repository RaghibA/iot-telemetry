package models

import "gorm.io/gorm"

type KafkaACL struct {
	gorm.Model
	ID          string
	APIKey      string
	TopicName   string
	Permissions []string
}
