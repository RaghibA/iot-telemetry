package app

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/iot-admin/internal/db"
	"github.com/RaghibA/iot-telemetry/iot-admin/internal/routes"
	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("Starting admin service")

	db.Connect()
	db.DeviceMigrate()
	db.BucketMigrate()

	r := gin.Default()

	routes.Admin(r)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
