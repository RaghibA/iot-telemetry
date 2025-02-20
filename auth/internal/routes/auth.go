package routes

import (
	"github.com/RaghibA/iot-telemetry/auth-service/internal/handlers"
	"github.com/RaghibA/iot-telemetry/auth-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

/**
 * Auth router group
 *
 * /auth/health: Service Health Check
 * /auth/register: Account Creation
 * /auth/login: account login
 * /auth/access-token: Access token generation
 * /auth/logout: user logout/cookie delete
 * /auth/deactivate: delete account and associated resources
 * /auth/api-ket: generate new api key
 *
 */
func Auth(r *gin.Engine) {
	userRouter := r.Group("/auth", middleware.Monitoring()) // monitoring middleware applied to all routes
	{
		userRouter.GET("/health", handlers.HealthHandler)
		userRouter.POST("/register", handlers.RegisterUserHandler)
		userRouter.POST("/login", handlers.LoginHandler)
		userRouter.POST("/access-token", middleware.Authenticate(), handlers.AccessTokenHandler)
		userRouter.POST("/logout", handlers.LogoutHandler)
		userRouter.DELETE("/deactivate", middleware.Authenticate(), handlers.DeactivateHandler)
		userRouter.GET("/api-key", middleware.Authenticate(), handlers.GenerateAPIKeyHandler)
	}
}
