package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
