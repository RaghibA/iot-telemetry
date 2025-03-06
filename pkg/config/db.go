package config

import (
    "fmt"

    "github.com/RaghibA/iot-telemetry/pkg/utils"
)

type DBConfig struct {
    PostgresUser string
    PostgresPass string
    PostgresName string
    PostgresPort string
}

// GetDBString constructs the database connection string from the given DBConfig.
// Params:
// - config: *DBConfig - a pointer to the DBConfig struct containing the database configuration
// Returns:
// - string: the constructed database connection string
func GetDBString(config *DBConfig) string {
    return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
        config.PostgresUser,
        config.PostgresPass,
        "db",
        config.PostgresPort,
        config.PostgresName,
    )
}

// GetDBConfig retrieves the database configuration from environment variables.
// Params: None
// Returns:
// - *DBConfig: a pointer to the DBConfig struct containing the database configuration
// - error: error if any occurred during the retrieval of environment variables
func GetDBConfig() (*DBConfig, error) {
    user, err := utils.GetEnv("DB_USER", "")
    if err != nil {
        return nil, err
    }

    pass, err := utils.GetEnv("DB_PASS", "")
    if err != nil {
        return nil, err
    }

    name, err := utils.GetEnv("DB_NAME", "")
    if err != nil {
        return nil, err
    }

    port, err := utils.GetEnv("DB_PORT", "")
    if err != nil {
        return nil, err
    }

    return &DBConfig{
        PostgresUser: user,
        PostgresPass: pass,
        PostgresName: name,
        PostgresPort: port,
    }, nil
}
