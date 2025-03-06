package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/pkg/jwt"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/consumer/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	store  store.ConsumerStore
	logger *log.Logger
	kafka  kafka.KafkaClient
}

func NewConsumerHander(store store.ConsumerStore, logger *log.Logger, kafka kafka.KafkaClient) *Handler {
	return &Handler{store: store, logger: logger, kafka: kafka}
}

func (h *Handler) ConsumerRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/messages", jwt.AuthWithAccessToken(h.ConsumerMessages)).Methods(http.MethodGet)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Health OK",
	})
}

func (h *Handler) ConsumerMessages(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Println("WS upgrade failed:", err)
		return
	}
	defer conn.Close()

	userId := r.Context().Value(jwt.UserKey).(string)

	if userId == "" {
		h.logger.Println("no user id found in ctx")
		return
	}

	deviceId := r.Header.Get("x-device-id")
	if deviceId == "" {
		h.logger.Println("no device id in req header")
		return
	}

	dbCtx := context.Background()
	device, err := h.store.GetDeviceById(dbCtx, deviceId)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.logger.Println("device not found: ", err)
		} else {
			h.logger.Println(err)
		}
		return
	}

	if userId != device.UserID {
		h.logger.Println("device uid & access token uid mismatch", userId, device.UserID)
		return
	}

	h.kafka.ConsumeFromTopic(device.TopicName, deviceId, conn)
}
