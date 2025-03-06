package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/services/auth/internal/monitoring"
	"github.com/RaghibA/iot-telemetry/services/auth/internal/routes"
	"github.com/RaghibA/iot-telemetry/services/auth/internal/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type AuthServer struct {
	Addr   string
	Db     *pgx.Conn
	Logger *log.Logger
}

// NewAuthServer creates a new authentication server instance.
// Params:
// - config: *config.AuthConfig - the authentication configuration
// - db: *pgx.Conn - the database connection
// Returns:
// - *AuthServer: a pointer to the created AuthServer
func NewAuthServer(config *config.AuthConfig, db *pgx.Conn) *AuthServer {
	addr := fmt.Sprintf("%s:%s", config.HOST, config.PORT)
	logger := log.New(os.Stdout, "AUTH_SERVER: ", log.LstdFlags)

	return &AuthServer{
		Addr:   addr,
		Db:     db,
		Logger: logger,
	}
}

// Run starts the authentication server.
// Params: None
// Returns:
// - error: error if any occurred during the server startup
func (s *AuthServer) Run() error {
	metrics := monitoring.NewMetrics()

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	subRouter.Use(metrics.MetricMonitoring) // Apply middleware to only this subrouter

	userStore := store.NewUserStore(s.Db, s.Logger)
	userHandler := routes.NewUserHandler(userStore, s.Logger)
	userHandler.UserRoutes(subRouter)

	router.Handle("/metrics", monitoring.PrometheusHandler()) // Expose metrics at /metrics

	log.Printf("Auth server running on %v", s.Addr)
	return http.ListenAndServe(s.Addr, router)
}
