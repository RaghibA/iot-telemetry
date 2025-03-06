package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/data/internal/monitoring"
	"github.com/RaghibA/iot-telemetry/services/data/internal/routes"
	"github.com/RaghibA/iot-telemetry/services/data/internal/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type TelemetryServer struct {
	addr        string
	db          *pgx.Conn
	logger      *log.Logger
	kafkaClient kafka.KafkaClient
}

func NewDataServer(config *config.DataConfig, db *pgx.Conn, logger *log.Logger, kafkaClient kafka.KafkaClient) *TelemetryServer {
	addr := fmt.Sprintf("%s:%s", config.HOST, config.PORT)
	return &TelemetryServer{
		addr:        addr,
		db:          db,
		logger:      logger,
		kafkaClient: kafkaClient,
	}
}

func (s *TelemetryServer) Run() error {
	metrics := monitoring.NewMetrics()

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1/data").Subrouter()
	subRouter.Use(metrics.MetricMonitoring) // Apply middleware to only this subrouter

	eventStore := store.NewEventStore(s.db, s.logger)
	dataHandler := routes.NewDataHandler(eventStore, s.logger, s.kafkaClient)
	dataHandler.DataRoutes(subRouter)

	router.Handle("/metrics", monitoring.PrometheusHandler()) // Expose metrics at /metrics

	log.Printf("Data server running on %v", s.addr)
	return http.ListenAndServe(s.addr, router)
}
