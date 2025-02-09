package app

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/db"
	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("starting data service")

	db.Connect()

	r := gin.Default()
	routes.Data(r)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
