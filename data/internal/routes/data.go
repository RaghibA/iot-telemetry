package routes

import (
	"github.com/RaghibA/iot-telemetry/iot-data-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

/**
 * Data Router Group
 *
 * GET: /telemetry/health -> Data service health check
 * POST: /telemetry/send -> Send device telemetry to kafka topic
 *
 */
func Data(r *gin.Engine) {
	dataRouter := r.Group("/telemetry")
	{
		dataRouter.GET("/health", handlers.HealthHandler)
		dataRouter.POST("/send", handlers.SendTelemetryHandler)
	}
}
