package models

import "time"

type ApiKey struct {
	UserID    string
	APIKey    string
	CreatedAt time.Time
}
