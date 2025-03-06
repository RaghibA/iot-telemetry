package app

import (
	"log"
	"os"

	"github.com/RaghibA/iot-telemetry/db"
	"github.com/RaghibA/iot-telemetry/pkg/config"
	"github.com/RaghibA/iot-telemetry/pkg/kafka"
	"github.com/RaghibA/iot-telemetry/services/consumer/internal/server"
)

func Run() {
	dbConfig, err := config.GetDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.NewDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	consumerConfig, err := config.GetConsumerConfig()
	if err != nil {
		log.Fatal(err)
	}

	kc := &kafka.KafkaService{}
	logger := log.New(os.Stdout, "CONSUMER SERVICE: ", log.LstdFlags)
	s := server.NewConsumerServer(consumerConfig, db, logger, kc)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
