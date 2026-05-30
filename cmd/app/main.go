package main

import (
	"log"

	"github.com/Jeno7u/server-app-course/internal/app"
	"github.com/Jeno7u/server-app-course/internal/db"
	"github.com/Jeno7u/server-app-course/internal/validation"
)

func main() {
	db.InitDB()
	validation.RegisterValidators()

	router := app.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
