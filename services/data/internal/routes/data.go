package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/data/internal/store"
	"github.com/gorilla/mux"
)

type Handler struct {
	store  store.EventStore
	logger *log.Logger
	kafka  kafka.KafkaClient
}

type SendEventRequestBody struct {
	DeviceID string          `json:"deviceId"`
	Data     json.RawMessage `json:"data"`
}

func NewDataHandler(store store.EventStore, logger *log.Logger, kafka kafka.KafkaClient) *Handler {
	return &Handler{store: store, logger: logger, kafka: kafka}
}

func (h *Handler) DataRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/event", h.sendTelemetry).Methods(http.MethodPost)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Health OK",
	})
}

func (h *Handler) sendTelemetry(w http.ResponseWriter, r *http.Request) {
	var eventData SendEventRequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&eventData); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deviceId := eventData.DeviceID
	if deviceId == "" {
		log.Println("No device id")
		http.Error(w, "Provide deviceId in request body", http.StatusBadRequest)
		return
	}

	apiKeyString := r.Header.Get("x-api-key")
	if apiKeyString == "" {
		h.logger.Println("No api key in header")
		http.Error(w, "Provide api key in 'x-api-key' header", http.StatusBadRequest)
		return
	}

	dbCtx := context.Background()
	apiKey, err := h.store.GetApiKey(dbCtx, apiKeyString)
	if err != nil {
		log.Println("db get api key", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	device, err := h.store.GetDeviceByDeviceId(dbCtx, deviceId)
	if err != nil {
		log.Println("db get device", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if apiKey.UserID != device.UserID {
		log.Println("api key & device userId missmatch")
		http.Error(w, "API key provided does not have permission to send data from this device", http.StatusUnauthorized)
		return
	}

	err = h.kafka.SendTelemetry(eventData.Data, device.TopicName, device.DeviceID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Telemetry Sent",
	})
}
