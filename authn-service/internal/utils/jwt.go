package utils

import (
	"errors"
	"os"
	"time"

	"github.com/RaghibA/iot-telemetry/authn-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(isAccessToken bool, account models.User, exp int64, iat int64) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"sub": account.UserID,
		"exp": exp,
		"iat": iat,
	}

	if isAccessToken {
		claims["permissions"] = []string{"read", "write"}
	} else {
		claims["permissions"] = []string{""}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func AuthenticateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // validate signing method
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	exp, ok := claims["exp"].(float64) // get exp from claims map / check it exists
	if ok {                            // if exists validate time
		expTime := time.Unix(int64(exp), 0)
		if expTime.Before(time.Now()) { // if time > now: expired
			return nil, errors.New("token expired")
		}
	} else { // if not ok, claim is missing
		return nil, errors.New("missing exp claim")
	}

	return claims, nil
}
