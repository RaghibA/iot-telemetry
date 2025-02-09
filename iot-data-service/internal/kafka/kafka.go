package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

func SendTelemetry(payload json.RawMessage, topic string, deviceID string) error {
	broker := []string{fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))}

	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		log.Println("failed to create producer", err)
		return err
	}

	defer producer.Close()

	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(deviceID),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := producer.SendMessage(&msg)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Message Sent to producer %d offset %d", partition, offset)

	return nil
}
