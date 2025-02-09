package routes

import (
	"github.com/RaghibA/iot-telemetry/iot-admin/internal/handlers"
	"github.com/RaghibA/iot-telemetry/iot-admin/internal/middleware"
	"github.com/gin-gonic/gin"
)

/**
 * Admin router group
 *
 * GET: /admin/health -> Admin healthcheck
 * POST: /admin/device -> Creates device & associated resources
 * GET: /admin/device -> Gets all user device info
 * GET: /admin/device?id=<deviceId> -> Gets single device info
 * DELETE: /admin/device -> Deletes device & associated resources
 *
 */
func Admin(r *gin.Engine) {
	adminRouter := r.Group("/admin")
	{
		adminRouter.GET("/health", handlers.HealthCheckHandler)
		adminRouter.POST("/device", middleware.Authenticate(), handlers.RegisterDeviceHandler)
		adminRouter.GET("/device", middleware.Authenticate(), handlers.GetDevicesHandler)
		adminRouter.DELETE("/device", middleware.Authenticate(), handlers.DeleteDeviceHandler)
	}
}
