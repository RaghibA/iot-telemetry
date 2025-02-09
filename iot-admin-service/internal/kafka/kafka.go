package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

/**
 * Generates a topic name based with 'topic.<deviceName>.<deviceId>.read' convention
 *
 * @params (deviceName, deviceId): ensures topic names are unique
 * @output string: single string for topic name
 */
func GenerateTopicName(deviceName string, deviceId string) string {
	return fmt.Sprintf("topic.%s.%s.read", strings.ReplaceAll(deviceName, " ", "-"), deviceId)
}

/**
 * Creates A topic in kafka
 *
 * @params topicName: topic name string defined by topic generation convention
 * @output error: error if topic generation fails
 */
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

/**
 * deletes device topic from kafka
 *
 * @params topicName: topic name used to remove topic from kafka server
 * @output error: error if deletion fails
 */
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
