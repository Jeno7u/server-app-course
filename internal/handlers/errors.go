package handlers

import (
	"github.com/Jeno7u/server-app-course/internal/errors"
	"github.com/gin-gonic/gin"
)

func TriggerErrorA(c *gin.Context) {
	_ = c.Error(&errors.CustomErrorA{Message: "Custom Error A triggered manually"})
}

func TriggerErrorB(c *gin.Context) {
	_ = c.Error(&errors.CustomErrorB{Message: "Custom Error B triggered manually (Not Found)"})
}
