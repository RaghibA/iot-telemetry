package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
)

// KafkaClient defines the interface for Kafka operations.
type KafkaClient interface {
	GenerateTopicName(deviceName string, deviceId string) string
	CreateTopic(topicName string) error
	DeleteTopic(topicName string) error
	SendTelemetry(payload json.RawMessage, topic string, deviceID string) error
	ConsumeFromTopic(topic string, deviceID string, conn *websocket.Conn)
}

type KafkaService struct{}

// GenerateTopicName generates a topic name based on the device name and device ID.
// Params:
// - deviceName: string - the name of the device
// - deviceId: string - the ID of the device
// Returns:
// - string: the generated topic name
func (k *KafkaService) GenerateTopicName(deviceName string, deviceId string) string {
	return fmt.Sprintf("topic.%s.%s.read", strings.ReplaceAll(deviceName, " ", "-"), deviceId)
}

// CreateTopic creates a new topic in Kafka.
// Params:
// - topicName: string - the name of the topic to create
// Returns:
// - error: error if any occurred during the topic creation
func (k *KafkaService) CreateTopic(topicName string) error {
	broker := []string{fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))}

	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0

	admin, err := sarama.NewClusterAdmin(broker, config)
	if err != nil {
		log.Println("error creating client:", err)
		return err
	}

	defer func() { _ = admin.Close() }()
	err = admin.CreateTopic(topicName, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteTopic deletes a topic in Kafka.
// Params:
// - topicName: string - the name of the topic to delete
// Returns:
// - error: error if any occurred during the topic deletion
func (k *KafkaService) DeleteTopic(topicName string) error {
	broker := []string{fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))}

	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0

	admin, err := sarama.NewClusterAdmin(broker, config)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() { _ = admin.Close() }()
	err = admin.DeleteTopic(topicName)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (k *KafkaService) SendTelemetry(payload json.RawMessage, topic string, deviceID string) error {
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

func (k *KafkaService) ConsumeFromTopic(topic string, deviceID string, conn *websocket.Conn) {
	broker := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))

	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log.Println("error creating kafka consumer", err)
		return
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Println("error getting partitions", err)
		return
	}

	pConsumer, err := consumer.ConsumePartition(topic, partitions[0], sarama.OffsetNewest)
	if err != nil {
		log.Println("failed to consume messages", err)
		return
	}
	defer pConsumer.Close()

	for message := range pConsumer.Messages() {
		err := conn.WriteMessage(websocket.TextMessage, message.Value)
		if err != nil {
			log.Println("failed to write message to ws writer", err)
			return
		}
	}
}
