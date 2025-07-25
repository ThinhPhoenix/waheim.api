package middleware

import (
	"github.com/gin-gonic/gin"
)

func RequireAuthorize(c *gin.Context) {
	if c.GetHeader("Authorization") == "" {
		cookie, err := c.Request.Cookie("token")
		if err == nil && cookie.Value != "" {
			c.Request.Header.Set("Authorization", "Bearer "+cookie.Value)
		}
	}
	c.Next()
}
