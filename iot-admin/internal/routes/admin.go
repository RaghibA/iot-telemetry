package routes

import (
	"github.com/RaghibA/iot-telemetry/iot-admin/internal/handlers"
	"github.com/RaghibA/iot-telemetry/iot-admin/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Admin(r *gin.Engine) {
	adminRouter := r.Group("/admin")
	{
		adminRouter.GET("/health", handlers.HealthCheckHandler)
		adminRouter.POST("/device", middleware.Authenticate(), handlers.RegisterDeviceHandler)
		adminRouter.GET("/device", middleware.Authenticate(), handlers.GetDevicesHandler)
		adminRouter.DELETE("/device", middleware.Authenticate(), handlers.DeleteDeviceHandler)
	}
}
