package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/pkg/jwt"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/consumer/internal/monitoring"
	"github.com/RaghibA/iot-telemetry/services/consumer/internal/routes"
	"github.com/RaghibA/iot-telemetry/services/consumer/internal/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type ConsumerServer struct {
	addr        string
	db          *pgx.Conn
	logger      *log.Logger
	kafkaClient kafka.KafkaClient
}

func NewConsumerServer(config *config.ConsumerConfig, db *pgx.Conn, logger *log.Logger, kafkaClient kafka.KafkaClient) *ConsumerServer {
	return &ConsumerServer{
		addr:        fmt.Sprintf("%s:%s", config.HOST, config.PORT),
		db:          db,
		logger:      logger,
		kafkaClient: kafkaClient,
	}
}

func (s *ConsumerServer) Run() error {
	metrics := monitoring.NewMetrics()

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1/telemetry").Subrouter()
	subRouter.Use(metrics.MetricMonitoring) // Apply middleware to only this subrouter

	consumerStore := store.NewConsumerStore(s.db, s.logger)
	consumerHandler := routes.NewConsumerHander(consumerStore, s.logger, s.kafkaClient)
	consumerHandler.ConsumerRoutes(subRouter)

	router.NewRoute().Path("/api/v1/telemetry/ws").HandlerFunc(jwt.AuthWithAccessToken(consumerHandler.ConsumerMessages))

	router.Handle("/metrics", monitoring.PrometheusHandler()) // Expose metrics at /metrics

	log.Printf("Consumer server running on %v", s.addr)
	return http.ListenAndServe(s.addr, router)
}
