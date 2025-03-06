package db

import (
	"context"
	"fmt"
	"log"

	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/jackc/pgx/v5"
)

// NewDB creates a new database connection.
// Params:
// - config: *config.DBConfig - the database configuration
// Returns:
// - *pgx.Conn: a pointer to the established database connection
// - error: error if any occurred during the connection establishment
func NewDB(config *config.DBConfig) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=UTC",
		config.PostgresUser,
		config.PostgresPass,
		config.PostgresName,
		config.PostgresPort,
	)
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	log.Println("DB Connection Established")

	return db, nil
}
