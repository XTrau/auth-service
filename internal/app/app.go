package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XTrau/auth-service/internal/database"
)

func Run() error {
	// Конфиг
	cfg := LoadConfig()

	// Подключение к бд
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к Postgres: %w", err)
	}

	// Регистрация хендлеров
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Запуск сервера
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Ошибка при работе сервера", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := server.Shutdown(stopCtx); err != nil {
		return fmt.Errorf("Ошибка при выключении сервера: %w", err)
	}

	return nil
}
