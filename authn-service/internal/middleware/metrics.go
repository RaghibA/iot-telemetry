package middleware

import (
	"strconv"
	"time"

	"github.com/RaghibA/iot-telemetry/authn-service/internal/monitoring"
	"github.com/gin-gonic/gin"
)

/**
 * Metric monitoring with prometheus
 *
 * @output - gin.HandlerFunc as middleware
 */
func Monitoring() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		dur := time.Since(start).Seconds()
		monitoring.HttpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(dur)

		statusCode := strconv.Itoa(c.Writer.Status())
		monitoring.HttpRequestStatus.WithLabelValues(c.Request.Method, c.FullPath(), statusCode).Inc()
	}
}
