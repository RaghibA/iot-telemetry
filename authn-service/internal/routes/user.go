package routes

import (
	"github.com/RaghibA/iot-telemetry/authn-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func User(r *gin.Engine) {
	userRouter := r.Group("/user")
	{
		userRouter.GET("/health", handlers.HealthHandler)
		userRouter.POST("/register", handlers.RegisterUserHandler)
		userRouter.POST("/login", handlers.LoginHandler)
		userRouter.POST("/refresh", handlers.RefreshHandler)
		userRouter.POST("/logout")
		userRouter.DELETE("/deactivate")
	}
}
