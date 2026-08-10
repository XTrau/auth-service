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
	"github.com/XTrau/auth-service/internal/repositories"
	"github.com/XTrau/auth-service/internal/usecases"
)

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
	}

	userRepository := repositories.NewPostgresUserRepository(db)

	jwtEncoder := authjwt.NewJwtEncoder(cfg)
	jwtDecoder := authjwt.NewJwtDecoder(cfg)
	jwtGenerator := authjwt.NewJwtGenerator(jwtEncoder)
	hasher := password.NewBcryptHasher(10)

	authUseCases := usecases.NewAuthUseCases(jwtGenerator, jwtDecoder, hasher, userRepository)

	// Регистрация хендлеров
	mux := http.NewServeMux()
	authHandlers := handlers.NewAuthHandlers(authUseCases)

	mux.HandleFunc("POST /auth/register", authHandlers.RegisterHandler)
	mux.HandleFunc("POST /auth/login", authHandlers.LoginHandler)
	mux.HandleFunc("POST /auth/logout", authHandlers.LogoutHandler)
	mux.HandleFunc("POST /auth/refresh", authHandlers.RefreshTokensHandler)
	mux.HandleFunc("GET /auth/user", authHandlers.GetCurrentUserHandler)

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
