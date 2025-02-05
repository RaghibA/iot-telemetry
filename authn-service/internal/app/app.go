package app

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/authn-service/internal/db"
	"github.com/RaghibA/iot-telemetry/authn-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("Running service")

	db.Connect()
	db.UserMigrate()

	r := gin.Default()

	// Register router groups
	routes.User(r)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
