package store

import (
	"context"
	"log"

	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/jackc/pgx/v5"
)

type UserStore interface {
	GetUserByID(ctx context.Context, userId string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	AddUser(ctx context.Context, user models.User) error
	AddApiKey(ctx context.Context, apikey models.ApiKey) error
	DeleteUser(ctx context.Context, userId string) error
	DeleteApiKey(ctx context.Context, userId string) error
}

type store struct {
	db     *pgx.Conn
	logger *log.Logger
}

// NewUserStore creates a new user store instance.
// Params:
// - db: *pgx.Conn - the database connection
// - logger: *log.Logger - the logger instance
// Returns:
// - *store: a pointer to the created store
func NewUserStore(db *pgx.Conn, logger *log.Logger) *store {
	return &store{db: db, logger: logger}
}

// GetUserByID retrieves a user by their ID.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user
// Returns:
// - *models.User: a pointer to the retrieved user
// - error: error if any occurred during the retrieval
func (s *store) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	var user models.User

	queryString := `
		SELECT user_id, username, password, email, created_at FROM users WHERE user_id=$1
	`
	err := s.db.QueryRow(ctx, queryString, userId).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
// Params:
// - ctx: context.Context - the context for the request
// - username: string - the username of the user
// Returns:
// - *models.User: a pointer to the retrieved user
// - error: error if any occurred during the retrieval
func (s *store) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	queryString := `
		SELECT user_id, username, password, email, created_at FROM users WHERE username=$1
	`
	err := s.db.QueryRow(ctx, queryString, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email.
// Params:
// - ctx: context.Context - the context for the request
// - email: string - the email of the user
// Returns:
// - *models.User: a pointer to the retrieved user
// - error: error if any occurred during the retrieval
func (s *store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	queryString := `
		SELECT user_id, username, password, email, created_at FROM users WHERE email=$1
	`
	err := s.db.QueryRow(ctx, queryString, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

// AddUser adds a new user to the database.
// Params:
// - ctx: context.Context - the context for the request
// - user: models.User - the user to add
// Returns:
// - error: error if any occurred during the addition
func (s *store) AddUser(ctx context.Context, user models.User) error {
	queryString := `
		INSERT INTO users (user_id, username, password, email)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(ctx, queryString, user.UserID, user.Username, user.Password, user.Email)
	return err
}

// AddApiKey adds a new API key to the database.
// Params:
// - ctx: context.Context - the context for the request
// - apikey: models.ApiKey - the API key to add
// Returns:
// - error: error if any occurred during the addition
func (s *store) AddApiKey(ctx context.Context, apikey models.ApiKey) error {
	queryString := `
		INSERT INTO api_keys (user_id, api_key)
		VALUES ($1, $2)
	`
	_, err := s.db.Exec(ctx, queryString, apikey.UserID, apikey.APIKey)
	return err
}

// DeleteUser deletes a user from the database.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user to delete
// Returns:
// - error: error if any occurred during the deletion
func (s *store) DeleteUser(ctx context.Context, userId string) error {
	queryString := `
		DELETE FROM users WHERE user_id=$1
	`

	_, err := s.db.Exec(ctx, queryString, userId)
	return err
}

// DeleteApiKey deletes an API key from the database.
// Params:
// - ctx: context.Context - the context for the request
// - userId: string - the ID of the user whose API key to delete
// Returns:
// - error: error if any occurred during the deletion
func (s *store) DeleteApiKey(ctx context.Context, userId string) error {
	queryString := `
		DELETE FROM api_keys WHERE user_id=$1
	`

	_, err := s.db.Exec(ctx, queryString, userId)
	return err
}
