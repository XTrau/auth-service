package main

import (
	"log"

	"github.com/XTrau/auth-service/internal/app"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic при работе сервиса: %v", err)
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервиса: %v", err.Error())
	}
}
