package store

import (
	"context"
	"log"

	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/jackc/pgx/v5"
)

// DeviceStore defines the interface for device-related database operations.
type DeviceStore interface {
	GetDeviceByID(ctx context.Context, deviceId string) (*models.Device, error)
	GetUserDevices(ctx context.Context, userId string) ([]models.Device, error)
	AddDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, deviceId string) error
}

type store struct {
	db     *pgx.Conn
	logger *log.Logger
}

// NewDeviceStore creates a new device store instance.
// Params:
// - db: *pgx.Conn - the database connection
// - logger: *log.Logger - the logger instance
// Returns:
// - *store: a pointer to the created store
func NewDeviceStore(db *pgx.Conn, logger *log.Logger) *store {
	return &store{db: db, logger: logger}
}

// GetDeviceByID retrieves a device by its ID.
// Params:
// - ctx: context.Context - the context for the request
// - deviceId: string - the ID of the device
// Returns:
// - *models.Device: a pointer to the retrieved device
// - error: error if any occurred during the retrieval
func (s *store) GetDeviceByID(ctx context.Context, deviceId string) (*models.Device, error) {
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

// GetUserDevices retrieves all devices belonging to a user.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user
// Returns:
// - []models.Device: slice of devices belonging to the user
// - error: error if any occurred during the retrieval
func (s *store) GetUserDevices(ctx context.Context, userId string) ([]models.Device, error) {
	var devices []models.Device

	queryString := `	
		SELECT device_name, device_id, user_id, topic_name, created_at FROM devices WHERE user_id=$1
	`

	rows, err := s.db.Query(ctx, queryString, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device models.Device

		err := rows.Scan(
			&device.DeviceName,
			&device.DeviceID,
			&device.UserID,
			&device.TopicName,
			&device.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// AddDevice adds a new device to the database.
// Params:
// - ctx: context.Context - the context for the request
// - device: *models.Device - pointer to the device to add
// Returns:
// - error: error if any occurred during the addition
func (s *store) AddDevice(ctx context.Context, device *models.Device) error {
	queryString := `
	INSERT INTO devices (device_name, device_id, user_id, topic_name)
	VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(
		ctx,
		queryString,
		device.DeviceName,
		device.DeviceID,
		device.UserID,
		device.TopicName,
	)
	return err
}

// DeleteDevice deletes a device from the database.
// Params:
// - ctx: context.Context - the context for the request
// - deviceId: string - the ID of the device to delete
// Returns:
// - error: error if any occurred during the deletion
func (s *store) DeleteDevice(ctx context.Context, deviceId string) error {
	queryString := `
    DELETE FROM devices WHERE device_id=$1
  `

	_, err := s.db.Exec(ctx, queryString, deviceId)
	return err
}
