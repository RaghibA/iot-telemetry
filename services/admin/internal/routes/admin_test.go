package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RaghibA/iot-telemetry/pkg/jwt"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/RaghibA/iot-telemetry/pkg/utils"
	"github.com/RaghibA/iot-telemetry/services/admin/internal/store"
	"github.com/gorilla/mux"
)

var testLogger *log.Logger
var buf *bytes.Buffer
var kc *kafka.MockKafkaServer

func TestMain(m *testing.M) {
	buf = new(bytes.Buffer)
	testLogger = utils.NewTestLogger(buf)
	kc = kafka.NewMockKafkaServer()

	code := m.Run()
	os.Exit(code)
}

func TestRegisterDeviceHandler(t *testing.T) {
	deviceStore := store.NewMockStore()
	handler := NewAdminHander(deviceStore, testLogger, kc)
	registerApi := "/api/v1/admin/device"

	t.Run("should fail if request body is invalid", func(t *testing.T) {
		buf.Reset()
		userId := "1234user"

		body := &CreateDeviceRequestBody{}
		marshalled, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, registerApi, bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(registerApi, jwt.AuthWithAccessToken(handler.registerDevice)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		if t.Failed() {
			t.Logf("Logs on failure: %s", buf.String())
		}
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		buf.Reset()

		userId := "test123"
		deviceName := "1234-test-device"

		body := &CreateDeviceRequestBody{
			DeviceName: deviceName,
		}
		marshalled, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(-time.Hour))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, registerApi, bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(registerApi, jwt.AuthWithAccessToken(handler.registerDevice)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should fail if device name exists in users devices", func(t *testing.T) {
		buf.Reset()

		userId := "test123"
		deviceName := "1234-test-device"

		deviceStore.Devices["sdkfksjlfkdjs"] = &models.Device{
			DeviceName: deviceName,
			UserID:     userId,
		}

		body := &CreateDeviceRequestBody{
			DeviceName: deviceName,
		}
		marshalled, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, registerApi, bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(registerApi, jwt.AuthWithAccessToken(handler.registerDevice)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("expected status code %d, got %d", http.StatusConflict, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should register a device", func(t *testing.T) {
		buf.Reset()
		for k := range deviceStore.Devices {
			delete(deviceStore.Devices, k)
		}

		userId := "test123"
		deviceName := "1234-test-device"

		body := &CreateDeviceRequestBody{
			DeviceName: deviceName,
		}
		marshalled, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, registerApi, bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(registerApi, jwt.AuthWithAccessToken(handler.registerDevice)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})
}

func TestGetDevicesHandler(t *testing.T) {
	getDevicesApi := "/api/v1/admin/device"
	deviceStore := store.NewMockStore()
	handler := NewAdminHander(deviceStore, testLogger, kc)
	userId := "1234test"
	deviceStore.AddDevice(context.Background(), &models.Device{
		DeviceName: "test1",
		DeviceID:   "test1234",
		UserID:     userId,
		TopicName:  kc.GenerateTopicName("test1", "test1234"),
		CreatedAt:  time.Now(),
	})
	deviceStore.AddDevice(context.Background(), &models.Device{
		DeviceName: "test2",
		DeviceID:   "test2345",
		UserID:     userId,
		TopicName:  kc.GenerateTopicName("test2", "test2345"),
		CreatedAt:  time.Now(),
	})

	t.Run("should fail if no access token is provided", func(t *testing.T) {
		buf.Reset()

		req, err := http.NewRequest(http.MethodGet, getDevicesApi, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(getDevicesApi, jwt.AuthWithAccessToken(handler.getDevices)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Logf("Logs on failure: %s", buf.String())
		}
	})

	t.Run("should return 400 if user has no devices", func(t *testing.T) {
		buf.Reset()
		newUserId := "4321user"

		token, err := jwt.GenerateAccessToken(newUserId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, getDevicesApi, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(getDevicesApi, jwt.AuthWithAccessToken(handler.getDevices)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		if t.Failed() {
			t.Logf("logs on failure: %s", buf.String())
		}
	})

	t.Run("should return 500 if db error occurs", func(t *testing.T) {
		buf.Reset()

		deviceStore.Err = errors.New("test error")
		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, getDevicesApi, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(getDevicesApi, jwt.AuthWithAccessToken(handler.getDevices)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		if t.Failed() {
			t.Logf("Logs on failure: %s", buf.String())
		}
	})

	t.Run("should get user devices", func(t *testing.T) {
		buf.Reset()

		deviceStore.Err = nil
		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, getDevicesApi, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(getDevicesApi, jwt.AuthWithAccessToken(handler.getDevices)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Logf("Logs on failure: %s", buf.String())
		}
	})
}

func TestDeleteDeviceHandler(t *testing.T) {
	deleteDeviceApi := "/api/v1/admin/device"
	deviceStore := store.NewMockStore()
	handler := NewAdminHander(deviceStore, testLogger, kc)
	userId := "1234test"
	deviceStore.AddDevice(context.Background(), &models.Device{
		DeviceName: "test1",
		DeviceID:   "test1234",
		UserID:     userId,
		TopicName:  kc.GenerateTopicName("test1", "test1234"),
		CreatedAt:  time.Now(),
	})

	t.Run("should fail if no device id is provided", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodDelete, deleteDeviceApi, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(deleteDeviceApi, jwt.AuthWithAccessToken(handler.deleteDevice)).Methods(http.MethodDelete)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		if t.Failed() {
			t.Log(buf.String())
		}
	})

	t.Run("should fail if no access token in provided", func(t *testing.T) {
		buf.Reset()

		req, err := http.NewRequest(http.MethodDelete, deleteDeviceApi, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(deleteDeviceApi, jwt.AuthWithAccessToken(handler.deleteDevice)).Methods(http.MethodDelete)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Log(buf.String())
		}
	})

	t.Run("should fail if access token user id does not match device user id", func(t *testing.T) {
		buf.Reset()

		newUserId := "32143132"
		token, err := jwt.GenerateAccessToken(newUserId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		wQueryParam := fmt.Sprintf("%s?deviceId=%s", deleteDeviceApi, "test1234")
		req, err := http.NewRequest(http.MethodDelete, wQueryParam, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(deleteDeviceApi, jwt.AuthWithAccessToken(handler.deleteDevice)).Methods(http.MethodDelete)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		if t.Failed() {
			t.Log(buf.String())
		}
	})

	t.Run("should delete user device", func(t *testing.T) {
		buf.Reset()

		token, err := jwt.GenerateAccessToken(userId, time.Now().Add(time.Hour*1))
		if err != nil {
			t.Fatal(err)
		}

		wQueryParam := fmt.Sprintf("%s?deviceId=%s", deleteDeviceApi, "test1234")
		req, err := http.NewRequest(http.MethodDelete, wQueryParam, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc(deleteDeviceApi, jwt.AuthWithAccessToken(handler.deleteDevice)).Methods(http.MethodDelete)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if t.Failed() {
			t.Log(buf.String())
		}
	})
}
