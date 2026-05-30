package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/gin-gonic/gin"
)

type calculatePayload struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

var feedbacks []models.Feedback

func Index(c *gin.Context) {
	data, err := os.ReadFile("internal/src/index.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func Calculate(c *gin.Context) {
	var payload calculatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("%v, %v", payload.Num1, payload.Num2)
	c.JSON(http.StatusOK, gin.H{"result": payload.Num1 + payload.Num2})
}

func GetUsers(c *gin.Context) {
	user := models.UsersResponse{Id: 1, Name: "Mironov Boris"}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
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
}

func Feedback(c *gin.Context) {
	var payload models.Feedback
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedbacks = append(feedbacks, payload)
	c.JSON(http.StatusAccepted, gin.H{"message": "Feedback received. Thank you, " + payload.Name})
}

func Headers(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")
	acceptLang := c.GetHeader("Accept-Language")

	c.JSON(http.StatusOK, gin.H{
		"User-agent":      userAgent,
		"Accept-Language": acceptLang,
	})
}

func validateAcceptLanguageFormat(lang string) bool {
	pattern := `^([a-zA-Z*]+(?:-[a-zA-Z*]+)?(?:;q=\d+(?:\.\d+)?)?)(,\s*([a-zA-Z*]+(?:-[a-zA-Z*]+)?(?:;q=\d+(?:\.\d+)?)?))*$`
	matched, _ := regexp.MatchString(pattern, lang)
	return matched
}

func parseHeaders(c *gin.Context) (*models.CommonHeaders, bool) {
	var headers models.CommonHeaders
	if err := c.ShouldBindHeader(&headers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required headers: User-Agent, Accept-Language"})
		return nil, false
	}

	if !validateAcceptLanguageFormat(headers.AcceptLanguage) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Accept-Language format"})
		return nil, false
	}

	return &headers, true
}

func Info(c *gin.Context) {
	headers, ok := parseHeaders(c)
	if !ok {
		return
	}

	serverTime := time.Now().UTC().Format(time.RFC3339)
	c.Header("X-Server-Time", serverTime)

	c.JSON(http.StatusOK, gin.H{
		"message": "Добро пожаловать! Ваши заголовки успешно обработаны.",
		"headers": gin.H{
			"User-Agent":      headers.UserAgent,
			"Accept-Language": headers.AcceptLanguage,
		},
	})
}

func Docs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to API docs"})
}

func CreateUser2(c *gin.Context) {
	var payload models.UserCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payload)
}

var products = []models.Product{
	{ProductID: 123, Name: "Smartphone", Category: "Electronics", Price: 599.99},
	{ProductID: 456, Name: "Phone Case", Category: "Accessories", Price: 19.99},
	{ProductID: 789, Name: "Iphone", Category: "Electronics", Price: 1299.99},
	{ProductID: 101, Name: "Headphones", Category: "Accessories", Price: 99.99},
	{ProductID: 202, Name: "Smartwatch", Category: "Electronics", Price: 299.99},
}

func GetProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	for i := range products {
		if products[i].ProductID == productID {
			c.JSON(http.StatusOK, products[i])
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	category := c.Query("category")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
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
}
