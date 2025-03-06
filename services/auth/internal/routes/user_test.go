package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RaghibA/iot-telemetry/pkg/jwt"
	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/RaghibA/iot-telemetry/pkg/utils"
	"github.com/RaghibA/iot-telemetry/services/auth/internal/store"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var testLogger *log.Logger
var buf *bytes.Buffer

func TestMain(m *testing.M) {
	buf = new(bytes.Buffer)
	testLogger = utils.NewTestLogger(buf)
	code := m.Run()

	os.Exit(code)
}

func TestRegisterUserHandler(t *testing.T) {
	userStore := store.NewMockUserStore()
	handler := NewUserHandler(userStore, testLogger)

	userStore.Users["123"] = &models.User{
		UserID:    "123",
		Username:  "test-user",
		Password:  []byte("1234"),
		Email:     "test@gmail.com",
		CreatedAt: time.Now(),
	}

	t.Run("should fail if the username is not unique", func(t *testing.T) {
		buf.Reset()
		payload := models.User{
			UserID:    "123",
			Username:  "test-user", // dupe username
			Password:  []byte("1234"),
			Email:     "test@gmail.com",
			CreatedAt: time.Now(),
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/register", handler.createUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("expected status code %d, got %d", http.StatusConflict, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should fail if the email is not unique", func(t *testing.T) {
		buf.Reset()
		payload := models.User{
			UserID:    "123",
			Username:  "test-user1",
			Password:  []byte("1234"),
			Email:     "test@gmail.com", // dupe email
			CreatedAt: time.Now(),
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/register", handler.createUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("expected status code %d, got %d", http.StatusConflict, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should create a new account", func(t *testing.T) {
		buf.Reset()
		payload := models.User{
			UserID:    "12345678",
			Username:  "newuser",
			Password:  []byte("1234test"),
			Email:     "newuser@gmail.com",
			CreatedAt: time.Now(),
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/register", handler.createUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
}

func TestLoginUserHandler(t *testing.T) {
	passString := "mypassword"
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(passString), 14)
	if err != nil {
		t.Fatal(err)
	}

	userStore := &store.MockUserStore{
		Users: map[string]*models.User{
			"123": {
				UserID:    "123",
				Username:  "test-user",
				Password:  hashedPass,
				Email:     "test@gmail.com",
				CreatedAt: time.Now(),
			},
		},
	}
	handler := NewUserHandler(userStore, testLogger)

	t.Run("should fail if the password is incorrect", func(t *testing.T) {
		buf.Reset()
		payload := LoginRequestBody{
			Username: "test-user",
			Password: "wrongpass",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/login", handler.loginUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should fail if the password is empty", func(t *testing.T) {
		buf.Reset()
		payload := LoginRequestBody{
			Username: "test-user",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/login", handler.loginUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should fail if the username is empty", func(t *testing.T) {
		buf.Reset()
		payload := LoginRequestBody{
			Password: "wrongpass",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/login", handler.loginUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should log in", func(t *testing.T) {
		buf.Reset()
		payload := LoginRequestBody{
			Username: "test-user",
			Password: "mypassword",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/login", handler.loginUser).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("expected status code %d, got %d", http.StatusAccepted, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
}

func TestAccessTokenHandler(t *testing.T) {
	userStore := store.NewMockUserStore()
	handler := NewUserHandler(userStore, testLogger)
	userId := "1234user"

	t.Run("should fail if cookie is expired", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateCookie(userId, time.Now().Add(-time.Hour))
		if err != nil {
			t.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/access-token", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/access-token", jwt.AuthWithCookie(handler.generateToken)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
	t.Run("should grant access token", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateCookie(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/access-token", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/access-token", jwt.AuthWithCookie(handler.generateToken)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
}

func TestLogoutHandler(t *testing.T) {
	userStore := store.NewMockUserStore()
	handler := NewUserHandler(userStore, testLogger)
	userId := "1234user"

	t.Run("should fail if cookie is expired", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateCookie(userId, time.Now().Add(-time.Hour))
		if err != nil {
			t.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookie)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/logout", jwt.AuthWithCookie(handler.logoutUser)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should logout user", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateCookie(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookie)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/api/v1/auth/logout", jwt.AuthWithCookie(handler.logoutUser)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
}

func TestApiKeyHandler(t *testing.T) {
	userStore := store.NewMockUserStore()
	handler := NewUserHandler(userStore, testLogger)

	userStore.ApiKeys["user1"] = models.ApiKey{
		APIKey:    "1234",
		UserID:    "user1",
		CreatedAt: time.Now(),
	}

	t.Run("should generate new API key", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateCookie("user1", time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Secure:   false,
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/api-key", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookie)

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/auth/api-key", jwt.AuthWithCookie(handler.regenerateApiKey)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
}
