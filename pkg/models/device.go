package models

import "time"

type Device struct {
	DeviceName string
	DeviceID   string
	UserID     string
	TopicName  string
	CreatedAt  time.Time
}
