package main

import (
	"log"

	"github.com/XTrau/auth-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Error on app running: %v", err.Error())
	}
}
