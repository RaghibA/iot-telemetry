package handlers

import (
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/consumer-service/internal/clients"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/db"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/kafka"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/middleware"
	"github.com/RaghibA/iot-telemetry/consumer-service/internal/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConsumerHandler(cd *clients.ClientDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WS upgrade failed", err)
			return
		}
		defer conn.Close()

		log.Println("connection established", r.RemoteAddr)

		deviceID := r.Header.Get("x-device-id")
		if deviceID == "" {
			log.Println("Missing deviceID")
			return
		}

		userID, ok := r.Context().Value(middleware.UserKey{}).(string)
		if !ok {
			log.Println("missing userId in ctx")
			return
		}

		var device models.Device
		err = db.IotDb.Db.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&device).Error
		if err != nil {
			log.Println("unable to find device in db")
			return
		}
		topicName := device.TopicName
		log.Println(device.DeviceID, topicName)

		// Add client to dir
		cd.AddClient(deviceID, conn)
		log.Printf("client dir: added %s (deviceId=%s)", r.RemoteAddr, deviceID)

		// Consumer messages from kafka topic
		kafka.ConsumeFromTopic(topicName, deviceID, conn)

		defer func() {
			cd.RemoveClient(deviceID)
			log.Printf("client dir: removed %s (deviceID=%s)", r.RemoteAddr, deviceID)
		}()
	}
}
