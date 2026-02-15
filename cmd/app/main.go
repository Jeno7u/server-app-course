package main

import (
	"log"
	"net/http"
	"os"

	models "github.com/Jeno7u/server-app-course/internal/models"
	validation "github.com/Jeno7u/server-app-course/internal/validation"
	"github.com/gin-gonic/gin"
)

type calculatePayload struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

var feedbacks []models.Feedback

func main() {
	router := gin.Default()
	validation.RegisterValidators()

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
		user := models.UsersResponse{Id: 1, Name: "Mironov Boris"}
		c.JSON(http.StatusOK, user)
	})

	router.POST("/user", func(c *gin.Context) {
		var payload models.UserPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.UserResponse{
			Name: payload.Name,
			Age: payload.Age,
			IsAdult: payload.Age >= 18,
		}	
		c.JSON(http.StatusOK, user)
	})

	router.POST("/feedback", func(c *gin.Context) {
		var payload models.Feedback
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		feedbacks = append(feedbacks, payload)
		c.JSON(http.StatusAccepted, gin.H{"message": "Feedback received. Thank you, " + payload.Name})
	})


	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}