package routes

import (
	"github.com/RaghibA/iot-telemetry/authn-service/internal/handlers"
	"github.com/RaghibA/iot-telemetry/authn-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Auth(r *gin.Engine) {
	userRouter := r.Group("/auth")
	{
		userRouter.GET("/health", handlers.HealthHandler)
		userRouter.POST("/register", handlers.RegisterUserHandler)
		userRouter.POST("/login", handlers.LoginHandler)
		userRouter.POST("/access-token", middleware.Authenticate(), handlers.RefreshHandler)
		userRouter.POST("/logout", handlers.LogoutHandler)
		userRouter.DELETE("/deactivate", middleware.Authenticate(), handlers.DeactivateHandler)
	}
}
