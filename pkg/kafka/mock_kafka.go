package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

/*
	GenerateTopicName(deviceName string, deviceId string) string
	CreateTopic(topicName string) error
	DeleteTopic(topicName string) error
*/

type MockKafkaServer struct {
	Topics   map[string]bool
	Err      error
	Messages map[string]json.RawMessage
}

func NewMockKafkaServer() *MockKafkaServer {
	return &MockKafkaServer{
		Topics:   make(map[string]bool),
		Err:      nil,
		Messages: make(map[string]json.RawMessage),
	}
}

func (k *MockKafkaServer) GenerateTopicName(deviceName string, deviceId string) string {
	return fmt.Sprintf("testtopic-%s-%s", deviceName, deviceId)
}

func (k *MockKafkaServer) CreateTopic(topicName string) error {
	if k.Err != nil {
		return k.Err
	}

	k.Topics[topicName] = true
	return nil
}

func (k *MockKafkaServer) DeleteTopic(topicName string) error {
	if k.Err != nil {
		return k.Err
	}

	delete(k.Topics, topicName)
	return nil
}

func (k *MockKafkaServer) SendTelemetry(payload json.RawMessage, topic string, deviceID string) error {
	if k.Err != nil {
		return k.Err
	}

	k.Messages[topic] = payload
	return nil
}

func (k *MockKafkaServer) ConsumeFromTopic(topic string, deviceID string, conn *websocket.Conn) {
	_ = conn.WriteMessage(websocket.TextMessage, []byte("test"))
}
