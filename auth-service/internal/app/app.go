package app

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/auth-service/internal/db"
	"github.com/RaghibA/iot-telemetry/auth-service/internal/monitoring"
	"github.com/RaghibA/iot-telemetry/auth-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("Starting auth service")

	// Init Prometheus monitoring
	monitoring.InitPrometheus()

	// Start DB connection & auto-migration
	db.Connect()
	db.UserMigrate()
	db.ACLMigrate()

	// Init gin & register auth routes
	r := gin.Default()
	routes.Auth(r)
	r.GET("/auth/metrics", gin.WrapH(monitoring.PrometheusHandler()))

	// Start listening
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
