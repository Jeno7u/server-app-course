package main

import (
	"log"
	"net/http"
	"os"

	user "github.com/Jeno7u/server-app-course/internal/models"
	"github.com/gin-gonic/gin"
)

type calculatePayload struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		data, err := os.ReadFile("internal/src/index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	router.POST("/calculate", func(c *gin.Context) {
		var payload calculatePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("%v, %v", payload.Num1, payload.Num2)
		c.JSON(http.StatusOK, gin.H{"result": payload.Num1 + payload.Num2})
	})

	router.GET("/users", func(c *gin.Context) {
		user := user.User{Id: 1, Name: "Mironov Boris"}
		c.JSON(http.StatusOK, user)
	})


	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}