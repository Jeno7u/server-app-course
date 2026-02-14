package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Добро пожаловать в моё приложение Gin!",
    })
  })
  router.Run() // listens on 0.0.0.0:8080 by default
}