package main

import (
	"log"

	"github.com/XTrau/auth-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервиса: %v", err.Error())
	}
}
