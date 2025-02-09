package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/iot-admin-service-service/internal/db"
	"github.com/RaghibA/iot-telemetry/iot-admin-service-service/internal/kafka"
	"github.com/RaghibA/iot-telemetry/iot-admin-service-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

/**
 * Admin handler Res/Req body structs
 */
type DeviceBody struct {
	DeviceName string `json:"deviceName"`
}

type DeviceRes struct {
	DeviceName string `json:"deviceName"`
	DeviceID   string `json:"deviceID"`
}

/**
 * Admin service health check
 */
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Health OK",
	})
}

/**
 * Registers devices & stores them in db
 *
 * Associates device to user api key to create ACLs for each topic
 *
 * Creates new kafka topic for the device
 */
func RegisterDeviceHandler(c *gin.Context) {
	var requestBody DeviceBody
	validate := validator.New()

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Internal Server Error",
			"code":  500001,
		})
		return
	}

	err := validate.Var(requestBody.DeviceName, "required,min=6")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device name must be at least 6 charachters",
			"code":  400001,
		})
		return
	}

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"code":    500002,
		})
		return
	}

	var ct int64
	err = db.IotDb.Db.Model(&models.Device{}).Where("user_id = ? AND device_name = ?", userId, requestBody.DeviceName).Count(&ct).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Errro",
			"code":  500004,
		})
		return
	}

	var user models.User
	err = db.IotDb.Db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Errro",
			"code":  500004,
		})
		return
	}

	if ct > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You already have a device with this name",
			"code":    400001,
		})
		return
	}

	deviceName := requestBody.DeviceName
	deviceId := uuid.New().String()
	topicName := kafka.GenerateTopicName(deviceName, deviceId)

	err = kafka.CreateTopic(topicName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"code":    500004,
		})
		return
	}

	device := models.Device{
		DeviceName: deviceName,
		DeviceID:   deviceId,
		UserID:     userId.(string),
		TopicName:  topicName,
	}

	acl := models.KafkaACL{
		UserID:    userId.(string),
		DeviceID:  deviceId,
		APIKey:    user.APIKey,
		TopicName: topicName,
		Write:     true,
		Read:      true,
	}

	err = db.IotDb.Db.Create(&device).Error
	if err != nil {
		log.Println(err)
		_ = kafka.DeleteTopic(topicName) // roll back topic if db write fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500003,
		})
		return
	}

	err = db.IotDb.Db.Create(&acl).Error
	if err != nil {
		log.Println(err)
		_ = kafka.DeleteTopic(topicName) // roll back topic if db write fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500005,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Device Registered",
		"deviceName": device.DeviceName,
		"deviceID":   device.DeviceID,
		"topicName":  device.TopicName,
	})
}

/**
 * Gets all user devices, or if 'id' query parameter is sent,
 * gets single device
 */
func GetDevicesHandler(c *gin.Context) {
	deviceId := c.DefaultQuery("id", "")
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"code":    500001,
		})
		return
	}
	var devices []models.Device
	var err error

	if deviceId == "" {
		err = db.IotDb.Db.Where("user_id = ?", userId).Find(&devices).Error
	} else {
		err = db.IotDb.Db.Where("user_id = ? AND device_id = ?", userId, deviceId).Find(&devices).Error
		if len(devices) == 0 {
			err = fmt.Errorf("no device with id:%s found", deviceId)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("No device with id:%s", deviceId),
				"code":    400001,
			})
			return
		}
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500002,
		})
		return
	}

	var deviceRes []DeviceRes
	for _, d := range devices {
		deviceRes = append(deviceRes, DeviceRes{
			DeviceName: d.DeviceName,
			DeviceID:   d.DeviceID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": deviceRes,
	})
}

/**
 * Deletes single device from 'id' query parameter
 *
 * deletes associated ACLs & kafka topic
 */
func DeleteDeviceHandler(c *gin.Context) {
	deviceId := c.DefaultQuery("id", "")
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"code":    500001,
		})
		return
	}

	var device models.Device
	err := db.IotDb.Db.Where("user_id = ? AND device_id = ?", userId, deviceId).First(&device).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Device not found",
			"code":    400001,
		})
		return
	}

	err = kafka.DeleteTopic(device.TopicName)
	if err != nil {
		log.Println("failed to delete topic:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500002,
		})
		return
	}

	err = db.IotDb.Db.Delete(&device).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500003,
		})
		return
	}

	var acl models.KafkaACL
	err = db.IotDb.Db.Where("user_id = ? AND device_id = ?", userId, deviceId).First(&acl).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500004,
		})
		return
	}

	err = db.IotDb.Db.Delete(&acl).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500005,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}
