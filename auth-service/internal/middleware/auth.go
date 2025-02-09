package middleware

import (
	"log"
	"net/http"

	"github.com/RaghibA/iot-telemetry/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
)

/**
 * Authenticates user via cookie
 *
 * @output - gin.HandlerFunc as middleware
 */
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Cookie not found",
				"code":    400006,
			})
			c.Abort()
			return
		}

		claims, err := utils.AuthenticateToken(token)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"code":    401002,
			})
			c.Abort()
			return
		}

		// Check sub claim
		userId, ok := claims["sub"].(string)
		if !ok || userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid userID",
				"code":    401003,
			})
			c.Abort()
			return
		}

		perms, ok := claims["permissions"].([]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "missing permissions",
				"code":    401004,
			})
			c.Abort()
			return
		}

		c.Set("userID", userId)
		c.Set("permissions", perms)

		c.Next()
	}
}
