package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RaghibA/iot-telemetry/consumer-service/internal/clients"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/db"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/handlers"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/middleware"
)

/**
 * Init app here
 */
func Run() {
	log.Println("Starting consumer service")

	db.Connect()

	clientDir := clients.InitClientDirectory()

	http.HandleFunc("/telemetry/consume", middleware.Authenticate(handlers.ConsumerHandler(clientDir)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")), nil))
}
