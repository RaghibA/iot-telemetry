package store

import (
	"context"
	"log"

	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/jackc/pgx/v5"
)

type ConsumerStore interface {
	GetDeviceById(ctx context.Context, deviceId string) (*models.Device, error)
}

type store struct {
	db     *pgx.Conn
	logger *log.Logger
}

func NewConsumerStore(db *pgx.Conn, logger *log.Logger) *store {
	return &store{
		db:     db,
		logger: logger,
	}
}

func (s *store) GetDeviceById(ctx context.Context, deviceId string) (*models.Device, error) {
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
