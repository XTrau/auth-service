package main

import (
	"log"
	"log/slog"

	"github.com/XTrau/auth-service/internal/config"
	"github.com/XTrau/auth-service/internal/database"
)

func main() {
	// Конфиг
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфига: %v", err.Error())
	}

	// Подключение к бд
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении к Postgres: %v", err.Error())
	}
	defer db.Close()

	// Миграции
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Ошибка при загрузке миграций: %v", err.Error())
	} else {
		slog.Info("Миграции загружены")
	}
}
