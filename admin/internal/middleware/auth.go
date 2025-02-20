package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/RaghibA/iot-telemetry/iot-admin-service/internal/utils"
	"github.com/gin-gonic/gin"
)

/**
 * Authenticates admin requests via access tokens
 *
 * @output - gin.HandlerFunc as middleware
 */
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Cookie not found",
				"code":    400006,
			})
			c.Abort()
			return
		}

		cClaims, err := utils.AuthenticateToken(cookie)
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
		cUserId, ok := cClaims["sub"].(string)
		if !ok || cUserId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid userID",
				"code":    401003,
			})
			c.Abort()
			return
		}

		accessToken := c.GetHeader("Authorization")
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		log.Println(accessToken)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "No access token",
				"code":    401001,
			})
			c.Abort()
			return
		}

		claims, err := utils.AuthenticateToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token expired",
				"code":    401002,
			})
			c.Abort()
			return
		}

		// check permissions claims
		reqPerms := map[string]bool{
			"read":  false,
			"write": false,
		}
		perms, ok := claims["permissions"].([]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "missing permissions",
				"code":    401003,
			})
			c.Abort()
			return
		}

		for _, p := range perms {
			if permStr, ok := p.(string); ok {
				if _, exists := reqPerms[permStr]; exists {
					reqPerms[permStr] = true
				}
			}
		}

		userId := claims["sub"].(string)

		if userId != cUserId { // check access token uid vs refresh token uid
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token Missmatch. Sign in and request a new access token",
				"code":    401005,
			})
			c.Abort()
			return
		}

		if reqPerms["read"] && reqPerms["write"] {
			c.Set("userID", userId)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "missing permissions",
				"code":    401004,
			})
			c.Abort()
			return
		}
	}
}
