package config

import "github.com/RaghibA/iot-telemetry/pkg/utils"

type AuthConfig struct {
	HOST      string
	PORT      string
	JWTSECRET string
}

type AdminConfig struct {
	HOST      string
	PORT      string
	JWTSECRET string
}

type DataConfig struct {
	HOST      string
	PORT      string
	JWTSECRET string
}

type ConsumerConfig struct {
	HOST      string
	PORT      string
	JWTSECRET string
}

// GetAuthConfig retrieves the authentication configuration from environment variables.
// Params: None
// Returns:
// - *AuthConfig: a pointer to the AuthConfig struct containing the configuration
// - error: error if any occurred during the retrieval of environment variables
func GetAuthConfig() (*AuthConfig, error) {
	host, err := utils.GetEnv("HOST", "")
	if err != nil {
		return nil, err
	}

	port, err := utils.GetEnv("PORT", "")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := utils.GetEnv("JWT_SECRET", "")
	if err != nil {
		return nil, err
	}

	return &AuthConfig{
		HOST:      host,
		PORT:      port,
		JWTSECRET: jwtSecret,
	}, nil
}

// GetAdminConfig retrieves the admin api configuration from environment variables.
// Params: None
// Returns:
// - *AdminConfig: a pointer to the AdminConfig struct containing the configuration
// - error: error if any occurred during the retrieval of environment variables
func GetAdminConfig() (*AdminConfig, error) {
	host, err := utils.GetEnv("HOST", "")
	if err != nil {
		return nil, err
	}

	port, err := utils.GetEnv("PORT", "")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := utils.GetEnv("JWT_SECRET", "")
	if err != nil {
		return nil, err
	}

	return &AdminConfig{
		HOST:      host,
		PORT:      port,
		JWTSECRET: jwtSecret,
	}, nil
}

// GetDataConfig retrieves the data api configuration from environment variables.
// Params: None
// Returns:
// - *DataConfig: a pointer to the DataConfig struct containing the configuration
// - error: error if any occurred during the retrieval of environment variables
func GetDataConfig() (*DataConfig, error) {
	host, err := utils.GetEnv("HOST", "")
	if err != nil {
		return nil, err
	}

	port, err := utils.GetEnv("PORT", "")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := utils.GetEnv("JWT_SECRET", "")
	if err != nil {
		return nil, err
	}

	return &DataConfig{
		HOST:      host,
		PORT:      port,
		JWTSECRET: jwtSecret,
	}, nil
}

func GetConsumerConfig() (*ConsumerConfig, error) {
	host, err := utils.GetEnv("HOST", "")
	if err != nil {
		return nil, err
	}

	port, err := utils.GetEnv("PORT", "")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := utils.GetEnv("JWT_SECRET", "")
	if err != nil {
		return nil, err
	}

	return &ConsumerConfig{
		HOST:      host,
		PORT:      port,
		JWTSECRET: jwtSecret,
	}, nil
}
