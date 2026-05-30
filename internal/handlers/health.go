package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"env":    env,
	})
}
