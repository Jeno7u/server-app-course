package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
			Name:    payload.Name,
			Age:     payload.Age,
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

	router.POST("/create_user", func(c *gin.Context) {
		var payload models.UserCreateRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, payload)
	})

	products := []models.Product{
		{
			ProductID: 123,
			Name:      "Smartphone",
			Category:  "Electronics",
			Price:     599.99,
		},
		{
			ProductID: 456,
			Name:      "Phone Case",
			Category:  "Accessories",
			Price:     19.99,
		},
		{
			ProductID: 789,
			Name:      "Iphone",
			Category:  "Electronics",
			Price:     1299.99,
		},
		{
			ProductID: 101,
			Name:      "Headphones",
			Category:  "Accessories",
			Price:     99.99,
		},
		{
			ProductID: 202,
			Name:      "Smartwatch",
			Category:  "Electronics",
			Price:     299.99,
		},
	}

	router.GET("/product/:productID", func(c *gin.Context) {
		productID, err := strconv.Atoi(c.Param("productID"))
		if err != nil {
			log.Println("Ошибка при получении productID из path: ", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		for i := range products {
			if products[i].ProductID == productID {
				c.JSON(http.StatusOK, products[i])
			}
		}

		c.Status(http.StatusNotFound)
	})

	router.GET("/products/search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		if keyword == "" {
			log.Println("keyword is empty in /products/search")
			c.Status(http.StatusBadRequest)
			return
		}

		category := c.Query("category")
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			log.Println("Ошибка при получении productID из path: ", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		filtered := products
		for _, p := range products {
			if (strings.Contains(strings.ToLower(p.Name), strings.ToLower(keyword))) && (category != "" && p.Category == category) {
				filtered = append(filtered, p)

				limit--
				if limit <= 0 {
					break
				}
			}
		}

		c.JSON(http.StatusOK, filtered)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
