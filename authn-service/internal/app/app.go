package app

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/authn-service/internal/db"
	"github.com/RaghibA/iot-telemetry/authn-service/internal/monitoring"
	"github.com/RaghibA/iot-telemetry/authn-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("Starting auth service")

	monitoring.InitPrometheus()

	db.Connect()
	db.UserMigrate()
	db.ACLMigrate()

	r := gin.Default()

	// Register router groups
	routes.Auth(r)

	r.GET("/auth/metrics", gin.WrapH(monitoring.PrometheusHandler()))

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
