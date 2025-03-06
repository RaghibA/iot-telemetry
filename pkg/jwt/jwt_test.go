package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"test": "ok",
	})
}

func TestGenerateCookie(t *testing.T) {

	t.Run("should return error if no userId is provided", func(t *testing.T) {
		_, err := GenerateCookie("", time.Now().Add(time.Hour*1))

		if err == nil {
			t.Error("expected error, got jwt token")
		}
	})

	t.Run("should return cookie with claims", func(t *testing.T) {
		tokenString, err := GenerateCookie("test", time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if _, exists := claims["sub"]; !exists {
			t.Error("no sub claim")
		}

		if _, exists := claims["iat"]; !exists {
			t.Error("no iat claim")
		}

		if _, exists := claims["exp"]; !exists {
			t.Error("no exp claim")
		}
	})
}

func TestGenerateAccessToken(t *testing.T) {

	t.Run("should return error if no user id is provided", func(t *testing.T) {
		_, err := GenerateAccessToken("", time.Now().Add(time.Hour*1))

		if err == nil {
			t.Error("expected error, got jwt token")
		}
	})

	t.Run("should return access token with claims", func(t *testing.T) {
		tokenString, err := GenerateAccessToken("test", time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if _, exists := claims["sub"]; !exists {
			t.Error("no sub claim")
		}

		if _, exists := claims["iat"]; !exists {
			t.Error("no iat claim")
		}

		if _, exists := claims["exp"]; !exists {
			t.Error("no exp claim")
		}
	})
}

func TestCookieMiddleware(t *testing.T) {

	t.Run("should fail if no cookie is provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.HandleFunc("/", AuthWithCookie(mockHandler)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %v, got %v", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should fail if cookie is expired", func(t *testing.T) {
		cookieString, err := GenerateCookie("1234", time.Now().Add(-time.Hour))
		if err != nil {
			t.Fatal(err)
		}
		cookie := http.Cookie{
			Name:     "refresh_token",
			Value:    cookieString,
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(&cookie)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/", AuthWithCookie(mockHandler)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %v, got %v", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should authorize request when cookie is provided", func(t *testing.T) {
		cookieString, err := GenerateCookie("1234", time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}
		cookie := http.Cookie{
			Name:     "refresh_token",
			Value:    cookieString,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(&cookie)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/", AuthWithCookie(mockHandler)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %v, got %v", http.StatusOK, rr.Code)
		}
	})
}

func TestAccessTokenMiddleware(t *testing.T) {

	t.Run("should fail if no token is provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router.HandleFunc("/", AuthWithAccessToken(mockHandler)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %v, got %v", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should fail if token is expired", func(t *testing.T) {
		tokenString, err := GenerateAccessToken("1234", time.Now().Add(-time.Hour))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/", AuthWithCookie(mockHandler)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %v, got %v", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should authorize request when token is provided", func(t *testing.T) {
		tokenString, err := GenerateAccessToken("1234", time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/", AuthWithAccessToken(mockHandler)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %v, got %v", http.StatusOK, rr.Code)
		}
	})
}
