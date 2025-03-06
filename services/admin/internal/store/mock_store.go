package store

import (
	"context"

	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/jackc/pgx/v5"
)

type MockStore struct {
	Devices map[string]*models.Device
	Err     error
}

func NewMockStore() *MockStore {
	return &MockStore{
		Devices: make(map[string]*models.Device),
		Err:     nil,
	}
}

/*
	GetDeviceByID(ctx context.Context, deviceId string) (*models.Device, error)
	GetUserDevices(ctx context.Context, userId string) ([]models.Device, error)
	AddDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, deviceId string) error
*/

func (s *MockStore) GetDeviceByID(ctx context.Context, deviceId string) (*models.Device, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	device, exists := s.Devices[deviceId]
	if !exists {
		return nil, pgx.ErrNoRows
	}
	return device, nil
}

func (s *MockStore) GetUserDevices(ctx context.Context, userId string) ([]models.Device, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	var devices []models.Device
	for _, device := range s.Devices {
		if device.UserID == userId {
			devices = append(devices, *device)
		}
	}

	return devices, nil
}

func (s *MockStore) AddDevice(ctx context.Context, device *models.Device) error {
	if s.Err != nil {
		return s.Err
	}

	s.Devices[device.DeviceID] = device
	return nil
}

func (s *MockStore) DeleteDevice(ctx context.Context, deviceId string) error {
	if s.Err != nil {
		return s.Err
	}

	delete(s.Devices, deviceId)
	return nil
}
