package utils

import (
	"os"

	"github.com/RaghibA/iot-telemetry/authn-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(isAccessToken bool, account models.User, exp int64, iat int64) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"sub": account.ID,
		"exp": exp,
		"iat": iat,
	}

	if isAccessToken {
		claims["permissions"] = []string{"read", "write"}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signed, nil
}
