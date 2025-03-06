package app

import (
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/db"
	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/admin/internal/server"
)

// Run initializes the configuration, database, Kafka client, and starts the admin server.
// Params: None
// Returns: None
func Run() {
	dbConfig, err := config.GetDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.NewDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	adminConfig, err := config.GetAdminConfig()
	if err != nil {
		log.Fatal(err)
	}

	kc := &kafka.KafkaService{}
	logger := log.New(os.Stdout, "ADMIN SERVICE: ", log.LstdFlags)
	s := server.NewAdminServer(adminConfig, db, logger, kc)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
