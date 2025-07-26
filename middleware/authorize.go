package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"waheim.api/configs"
)

func RequireAuthorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			cookie, err := c.Request.Cookie("token")
			if err == nil && cookie.Value != "" {
				c.Request.Header.Set("Authorization", "Bearer "+cookie.Value)
			}
		}
		token := c.GetHeader("Authorization")
		if token == "" {
			cookie, err := c.Request.Cookie("token")
			if err == nil && cookie.Value != "" {
				token = "Bearer " + cookie.Value
			}
		}
		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing or invalid token"})
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := configs.ValidateJwt(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		userId, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)
		ctx := context.WithValue(c.Request.Context(), "user_id", userId)
		ctx = context.WithValue(ctx, "role", role)
		c.Request = c.Request.WithContext(ctx)
		// Nếu truyền roles, kiểm tra quyền
		if len(roles) > 0 {
			allowed := false
			for _, r := range roles {
				if r == role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(403, gin.H{"error": "Permission denied"})
				return
			}
		}
		c.Next()
	}
}
