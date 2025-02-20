package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

/**
 * creates kafka producer and sends device telemetry data to topic
 *
 * @params (payload, topic, deviceId): data from req body, topic name for associated device
 * 	& deviceId for partioning key
 *
 * @output error: if at any point consumer fails abort and return error
 */
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
