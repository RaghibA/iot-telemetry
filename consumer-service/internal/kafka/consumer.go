package kafka

import (
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
)

/**
 * Consumes messages from device topic
 *
 * @params (topic, deviceID, conn): topicname for given deviceID,
 * 	ref to websocket conn. Used to send write kafka message directly
 * 	to client
 */
func ConsumeFromTopic(topic string, deviceID string, conn *websocket.Conn) {
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
