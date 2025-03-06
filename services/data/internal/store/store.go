package store

import (
	"context"
	"log"

	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/jackc/pgx/v5"
)

type EventStore interface {
	GetApiKey(ctx context.Context, key string) (*models.ApiKey, error)
	GetDeviceByDeviceId(ctx context.Context, deviceId string) (*models.Device, error)
}

type store struct {
	db     *pgx.Conn
	logger *log.Logger
}

func NewEventStore(db *pgx.Conn, logger *log.Logger) *store {
	return &store{db: db, logger: logger}
}

func (s *store) GetApiKey(ctx context.Context, key string) (*models.ApiKey, error) {
	var apiKey models.ApiKey

	queryString := `
		SELECT user_id, api_key, created_at FROM api_keys WHERE api_key=$1
	`
	err := s.db.QueryRow(ctx, queryString, key).Scan(
		&apiKey.UserID,
		&apiKey.APIKey,
		&apiKey.CreatedAt,
	)
	if err != nil {
		return &models.ApiKey{}, err
	}

	return &apiKey, nil
}

func (s *store) GetDeviceByDeviceId(ctx context.Context, deviceId string) (*models.Device, error) {
	var device models.Device

	queryString := `
		SELECT device_name, device_id, user_id, topic_name, created_at FROM devices WHERE device_id=$1
	`

	err := s.db.QueryRow(ctx, queryString, deviceId).Scan(
		&device.DeviceName,
		&device.DeviceID,
		&device.UserID,
		&device.TopicName,
		&device.CreatedAt,
	)
	if err != nil {
		return &models.Device{}, err
	}

	return &device, nil
}
