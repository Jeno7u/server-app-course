package validation

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)


func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	err := v.RegisterValidation("notcontains", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, s := range strings.Split(fl.Param(), ";") {
			if strings.Contains(strings.ToLower(value), strings.ToLower(s)) {
				return false
			}
		}
		return true
	})
	if err != nil {
		log.Fatalf("got error when tried to register validators, %v", err)
	}
}
}