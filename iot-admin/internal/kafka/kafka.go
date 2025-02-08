package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

func GenerateTopicName(deviceName string, deviceId string) string {
	return fmt.Sprintf("topic.%s.%s.read", strings.ReplaceAll(deviceName, " ", "-"), deviceId)
}

func CreateTopic(topicName string) error {
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

func DeleteTopic(topicName string) error {
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
