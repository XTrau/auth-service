package handlers

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type AuthorizationHandlers interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
	RefreshTokensHandler(w http.ResponseWriter, r *http.Request)
	GetCurrentUserHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterHandlers(mux *http.ServeMux, authHandlers AuthorizationHandlers) {
	mux.HandleFunc("POST /api/auth/register", authHandlers.RegisterHandler)
	mux.HandleFunc("POST /api/auth/login", authHandlers.LoginHandler)
	mux.HandleFunc("POST /api/auth/logout", authHandlers.LogoutHandler)
	mux.HandleFunc("POST /api/auth/refresh", authHandlers.RefreshTokensHandler)
	mux.HandleFunc("GET /api/auth/user", authHandlers.GetCurrentUserHandler)

	// Swagger документация
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
}
