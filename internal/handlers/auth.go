package handlers

import (
	"net/http"

	"github.com/XTrau/auth-service/internal/domain"
)

type AuthHandlers struct {
	userRepository domain.UserRepository
}

func NewAuthHandlers(userRepository domain.UserRepository) *AuthHandlers {
	return &AuthHandlers{userRepository: userRepository}
}

func (ah *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
}
