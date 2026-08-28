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

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/auth/password"
	"github.com/XTrau/auth-service/internal/database"
	"github.com/XTrau/auth-service/internal/handlers"
	"github.com/XTrau/auth-service/internal/middlewares"
	"github.com/XTrau/auth-service/internal/repositories"
	"github.com/XTrau/auth-service/internal/service"

	_ "github.com/XTrau/auth-service/docs"
)

// @title       Swagger auth-api
// @version     1.0
// @description Authorization API.
// @BasePath    /
func Run() error {
	// Конфиг
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Ошибка при загрзке конфига: %w", err)
	}

	// Подключение к бд
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к Postgres: %w", err)
	} else {
		slog.Debug("Подключение к базе данных успешно")
	}
	defer db.Close()

	// Миграции
	if err := RunMigrations(db); err != nil {
		return fmt.Errorf("Ошибка при загрузке миграций: %w", err)
	}
	slog.Info("Миграции загружены")

	// Зависимости
	jwtTokenizer := authjwt.NewRS256Tokenizer(cfg)
	argon2Hasher := password.NewArgon2Hasher(password.Argon2DefaultParams())

	unitOfWork := repositories.NewPgUnitOfWork(db)
	userRepository := repositories.NewPgUserRepository()
	blockedTokensRepository := repositories.NewPgBlockedTokensRepository()

	authorizationService := service.NewAuthenticationService(jwtTokenizer, argon2Hasher, unitOfWork, userRepository, blockedTokensRepository)

	// Регистрация хендлеров
	mux := http.NewServeMux()

	authHandlers := handlers.NewAuthenticationHandlers(authorizationService)

	handlers.RegisterHandlers(mux, authHandlers)

	// Регистрация Middlewares
	var h http.Handler = mux

	h = middlewares.Recover(h)
	h = middlewares.Log(h)

	// Параметры сервера
	serverAddr := ":8080"
	server := http.Server{
		Addr:    serverAddr,
		Handler: h,

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Запуск сервера
	go func() {
		slog.Info("Сервер запущен", slog.String("address", "http://localhost:8080"))
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
