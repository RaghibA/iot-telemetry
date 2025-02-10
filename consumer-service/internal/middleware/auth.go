package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RaghibA/iot-telemetry/consumer-service/internal/utils"
)

type UserKey struct{}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if accessToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate signature
		claims, err := utils.AuthenticateToken(accessToken)
		if err != nil {
			http.Error(w, "Failed to authenticate access token", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"]
		if !ok {
			http.Error(w, "missing claims", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, userID)
		r = r.WithContext(ctx)

		next(w, r.WithContext(ctx))
	}
}
