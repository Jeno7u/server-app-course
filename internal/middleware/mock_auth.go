package middleware

import (
	"net/http"
	"strconv"

	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/gin-gonic/gin"
)

func MockAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.GetHeader("X-User-Id")
		if userIdStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "X-User-Id header missing"})
			c.Abort()
			return
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "X-User-Id must be an integer"})
			c.Abort()
			return
		}

		role := c.GetHeader("X-User-Role")
		if role == "" {
			role = "user"
		}

		c.Set("user", models.CurrentUser{
			ID:   userId,
			Role: role,
		})
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
			c.Abort()
			return
		}

		user := val.(models.CurrentUser)
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Admin role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
