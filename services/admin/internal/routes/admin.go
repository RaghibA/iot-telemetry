package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/pkg/jwt"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/pkg/models"
	"github.com/RaghibA/iot-telemetry/services/admin/internal/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	store  store.DeviceStore
	logger *log.Logger
	kafka  kafka.KafkaClient
}

type CreateDeviceRequestBody struct {
	DeviceName string
}

// NewAdminHander creates a new handler for admin routes.
// Params:
// - store: store.DeviceStore - the device store instance
// - logger: *log.Logger - the logger instance
// - kafka: kafka.KafkaClient - the Kafka client instance
// Returns:
// - *Handler: a pointer to the created Handler
func NewAdminHander(store store.DeviceStore, logger *log.Logger, kafka kafka.KafkaClient) *Handler {
	return &Handler{store: store, logger: logger, kafka: kafka}
}

// AdminRoutes sets up the admin routes.
// Params:
// - router: *mux.Router - the router instance
// Returns: None
func (h *Handler) AdminRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/device", jwt.AuthWithAccessToken(h.registerDevice)).Methods(http.MethodPost)
	router.HandleFunc("/device", jwt.AuthWithAccessToken(h.getDevices)).Methods(http.MethodGet)
	router.HandleFunc("/device", jwt.AuthWithAccessToken(h.deleteDevice)).Methods(http.MethodDelete)
}

// healthCheck is a handler for the health check endpoint.
// Params:
// - w: http.ResponseWriter - the HTTP response writer
// - r: *http.Request - the HTTP request
// Returns: None
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Health OK",
	})
}

// registerDevice is a handler for registering a new device.
// Params:
// - w: http.ResponseWriter - the HTTP response writer
// - r: *http.Request - the HTTP request
// Returns: None
func (h *Handler) registerDevice(w http.ResponseWriter, r *http.Request) {
	var deviceBody CreateDeviceRequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&deviceBody); err != nil {
		h.logger.Println("Failed to unmarshal input")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if deviceBody.DeviceName == "" {
		h.logger.Println("Empty device name field")
		http.Error(w, "Provide a device name", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(jwt.UserKey).(string)
	if userId == "" {
		h.logger.Println("No userId in context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	deviceId := uuid.New().String()

	topicName := h.kafka.GenerateTopicName(deviceBody.DeviceName, deviceId)

	newDevice := &models.Device{
		DeviceName: deviceBody.DeviceName,
		DeviceID:   deviceId,
		UserID:     userId,
		TopicName:  topicName,
	}

	err := h.kafka.CreateTopic(topicName)
	if err != nil {
		h.logger.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dbCtx := context.Background()
	devices, err := h.store.GetUserDevices(dbCtx, userId)
	if err != nil && err != pgx.ErrNoRows {
		h.logger.Println(err)
		err = h.kafka.DeleteTopic(topicName)
		if err != nil {
			h.logger.Println("Failed to delete new topic after device db write err", err)
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, d := range devices {
		if d.DeviceName == newDevice.DeviceName {
			h.logger.Println("Duplicate device name")
			http.Error(
				w,
				fmt.Sprintf("You already have a device with the name: %s", newDevice.DeviceName),
				http.StatusConflict,
			)
			return
		}
	}

	dbCtx = context.Background()
	err = h.store.AddDevice(dbCtx, newDevice)
	if err != nil {
		h.logger.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deviceName": newDevice.DeviceName,
		"deviceId":   newDevice.DeviceID,
	})
}

// getDevices is a handler for retrieving all devices for a user.
// Params:
// - w: http.ResponseWriter - the HTTP response writer
// - r: *http.Request - the HTTP request
// Returns: None
func (h *Handler) getDevices(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.UserKey).(string)

	dbCtx := context.Background()
	devices, err := h.store.GetUserDevices(dbCtx, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.logger.Println("No device found for user")
			http.Error(w, "No device found for provided id", http.StatusBadRequest)
		} else {
			h.logger.Println("Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if len(devices) == 0 {
		h.logger.Println("user has no devices")
		http.Error(w, "No devices found", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"devices": devices,
	})
}

// deleteDevice is a handler for deleting a device.
// Params:
// - w: http.ResponseWriter - the HTTP response writer
// - r: *http.Request - the HTTP request
// Returns: None
func (h *Handler) deleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("deviceId")
	if deviceId == "" {
		h.logger.Println("no device id provided in query param")
		http.Error(w, "No deviceId provided in query param", http.StatusBadRequest)
		return
	}

	dbCtx := context.Background()
	device, err := h.store.GetDeviceByID(dbCtx, deviceId)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.logger.Printf("No device found for id: %s", deviceId)
			http.Error(w, "No device found for provided id", http.StatusBadRequest)
		} else {
			h.logger.Println("Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	deviceUserId := device.UserID
	userIdClaim := r.Context().Value(jwt.UserKey)
	if userIdClaim == nil {
		h.logger.Println("No userId claim")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if deviceUserId != userIdClaim {
		h.logger.Println("device user id & claim user is mismatch")
		h.logger.Println(deviceUserId, userIdClaim)
		http.Error(w, "You are not authorized to delete this device", http.StatusUnauthorized)
		return
	}

	err = h.store.DeleteDevice(dbCtx, deviceId)
	if err != nil {
		h.logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.kafka.DeleteTopic(device.TopicName)
	if err != nil {
		h.logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	deviceName := device.DeviceName

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("%s deleted", deviceName),
	})
}
