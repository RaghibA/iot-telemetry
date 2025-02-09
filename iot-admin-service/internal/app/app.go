package app

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/iot-admin-service-service/internal/db"
	"github.com/RaghibA/iot-telemetry/iot-admin-service-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("Starting admin service")

	// Connect DB & run migrations
	db.Connect()
	db.DeviceMigrate()
	db.ACLMigrate()

	// Init router & attach admin routes
	r := gin.Default()
	routes.Admin(r)

	// start admin listener
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
