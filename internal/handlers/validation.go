package handlers

import (
	"net/http"

	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/gin-gonic/gin"
)

func ValidateUser(c *gin.Context) {
	var req models.CustomUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	phone := "Unknown"
	if req.Phone != nil {
		phone = *req.Phone
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User is valid",
		"data": gin.H{
			"username": req.Username,
			"age":      req.Age,
			"email":    req.Email,
			"phone":    phone,
		},
	})
}
