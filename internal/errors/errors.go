package errors
package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomErrorA struct {
	Message string
}

func (e *CustomErrorA) Error() string {
	return e.Message
}

type CustomErrorB struct {
	Message string
}

func (e *CustomErrorB) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch err := e.Err.(type) {
				case *CustomErrorA:
					c.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
					return
				case *CustomErrorB:
					c.JSON(http.StatusNotFound, ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
					return
				case validator.ValidationErrors:
					var errors []string
					for _, f := range err {
						errors = append(errors, fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", f.Field(), f.Tag()))
					}
					c.JSON(http.StatusUnprocessableEntity, gin.H{"detail": errors})
					return
				default:
					c.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
					return
				}
			}
		}
	}
}
