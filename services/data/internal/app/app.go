package app

import (
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/db"
	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/data/internal/server"
)

// Run initializes the configuration, database, Kafka client, and starts the data server.
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

	dataConfig, err := config.GetDataConfig()
	if err != nil {
		log.Fatal(err)
	}

	kc := &kafka.KafkaService{}
	logger := log.New(os.Stdout, "DATA SERVICE: ", log.LstdFlags)
	s := server.NewDataServer(dataConfig, db, logger, kc)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
