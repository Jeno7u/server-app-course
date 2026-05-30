package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	user, ok := fetchCurrentUser(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserByID(c *gin.Context) {
	// Dummy implementation for structure
	id := c.Param("user_id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"role": "user",
	})
}
