package store

import (
	"context"
	"errors"

	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/jackc/pgx/v5"
)

type MockUserStore struct {
	Users   map[string]*models.User
	ApiKeys map[string]models.ApiKey
	Err     error
}

// NewMockUserStore creates a new mock user store instance.
// Params: None
// Returns:
// - *MockUserStore: a pointer to the created MockUserStore
func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		Users:   make(map[string]*models.User),
		ApiKeys: make(map[string]models.ApiKey),
		Err:     nil,
	}
}

// GetUserByID retrieves a user by their ID from the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user
// Returns:
// - *models.User: a pointer to the retrieved user
// - error: error if any occurred during the retrieval
func (m *MockUserStore) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	user, exists := m.Users[userId]
	if !exists {
		return nil, pgx.ErrNoRows
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username from the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - username: string - the username of the user
// Returns:
// - *models.User: a pointer to the retrieved user
// - error: error if any occurred during the retrieval
func (m *MockUserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, pgx.ErrNoRows
}

// GetUserByEmail retrieves a user by their email from the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - email: string - the email of the user
// Returns:
// - *models.User: a pointer to the retrieved user
// - error: error if any occurred during the retrieval
func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, pgx.ErrNoRows
}

// AddUser adds a new user to the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - user: models.User - the user to add
// Returns:
// - error: error if any occurred during the addition
func (m *MockUserStore) AddUser(ctx context.Context, user models.User) error {
	if m.Err != nil {
		return m.Err
	}

	for _, u := range m.Users {
		if u.Username == user.Username {
			return errors.New("username already exists")
		}
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	m.Users[user.UserID] = &user
	return nil
}

// AddApiKey adds a new API key to the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - apikey: models.ApiKey - the API key to add
// Returns:
// - error: error if any occurred during the addition
func (m *MockUserStore) AddApiKey(ctx context.Context, apikey models.ApiKey) error {
	if m.Err != nil {
		return m.Err
	}

	m.ApiKeys[apikey.UserID] = apikey
	return nil
}

// DeleteUser deletes a user from the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user to delete
// Returns:
// - error: error if any occurred during the deletion
func (m *MockUserStore) DeleteUser(ctx context.Context, userId string) error {
	if m.Err != nil {
		return m.Err
	}

	delete(m.Users, userId)
	return nil
}

// DeleteApiKey deletes an API key from the mock store.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user whose API key to delete
// Returns:
// - error: error if any occurred during the deletion
func (m *MockUserStore) DeleteApiKey(ctx context.Context, userId string) error {
	if m.Err != nil {
		return m.Err
	}

	delete(m.ApiKeys, userId)
	return nil
}
