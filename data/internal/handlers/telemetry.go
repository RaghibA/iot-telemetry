package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/db"
	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/kafka"
	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/models"
	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/utils"
	"github.com/gin-gonic/gin"
)

/**
 * Data service req/res body structs
 */
type DeviceTelemetryReq struct {
	DeviceID string          `json:"deviceID"`
	Data     json.RawMessage `json:"data"`
}

/**
 * data service health check handler
 */
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Health OK",
	})
}

/**
 * validates data & acls; sends data to device topic
 */
func SendTelemetryHandler(c *gin.Context) {
	var event DeviceTelemetryReq

	if c.Bind(&event) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Payload",
			"code":    400001,
		})
		return
	}

	apiKey := c.GetHeader("x-api-key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Send API Key as request header(x-api-key)",
			"code":    401001,
		})
		return
	}

	var acl models.KafkaACL
	err := db.IotDb.Db.Where("device_id = ?", event.DeviceID).First(&acl).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500002,
		})
		return
	}

	// Validate API Key
	if !utils.VerifyAPIKey(apiKey, acl.APIKey) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "API Invalid",
			"code":    401002,
		})
		return
	}

	err = kafka.SendTelemetry(event.Data, acl.TopicName, acl.DeviceID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"code":    500003,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "telemetry sent",
	})
}
