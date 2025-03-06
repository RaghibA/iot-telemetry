package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/admin/internal/monitoring"
	"github.com/RaghibA/iot-telemetry/services/admin/internal/routes"
	"github.com/RaghibA/iot-telemetry/services/admin/internal/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type AdminServer struct {
	addr        string
	db          *pgx.Conn
	logger      *log.Logger
	kafkaClient kafka.KafkaClient
}

// NewAdminServer creates a new authentication server instance.
// Params:
// - config: *config.AdminConfig - the authentication configuration
// - db: *pgx.Conn - the database connection
// - logger: *log.Logger - the logger instance
// - kafkaClient: kafka.KafkaClient - the Kafka client instance
// Returns:
// - *AdminServer: a pointer to the created AdminServer
func NewAdminServer(config *config.AdminConfig, db *pgx.Conn, logger *log.Logger, kafkaClient kafka.KafkaClient) *AdminServer {
	addr := fmt.Sprintf("%s:%s", config.HOST, config.PORT)
	return &AdminServer{
		addr:        addr,
		db:          db,
		logger:      logger,
		kafkaClient: kafkaClient,
	}
}

// Run starts the authentication server.
// Params: None
// Returns:
// - error: error if any occurred during the server startup
func (s *AdminServer) Run() error {
	metrics := monitoring.NewMetrics()

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1/admin").Subrouter()
	subRouter.Use(metrics.MetricMonitoring) // Apply middleware to only this subrouter

	deviceStore := store.NewDeviceStore(s.db, s.logger)
	deviceHandler := routes.NewAdminHander(deviceStore, s.logger, s.kafkaClient)
	deviceHandler.AdminRoutes(subRouter)

	router.Handle("/metrics", monitoring.PrometheusHandler()) // Expose metrics at /metrics

	log.Printf("Admin server running on %v", s.addr)
	return http.ListenAndServe(s.addr, router)
}
