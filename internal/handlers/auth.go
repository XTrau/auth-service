package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/XTrau/auth-service/internal/dto"
	"github.com/XTrau/auth-service/internal/usecases"
	"github.com/asaskevich/govalidator"
)

type AuthHandlers struct {
	authUseCases *usecases.AuthUseCases
}

func NewAuthHandlers(authUseCases *usecases.AuthUseCases) *AuthHandlers {
	return &AuthHandlers{authUseCases: authUseCases}
}

func (ah *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("плохое тело запроса на /register", slog.String("error", err.Error()))
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if ok, err := govalidator.ValidateStruct(req); !ok || err != nil {
		slog.Info("ошибка при валидации тела /register", slog.Any("тело", req), slog.String("error", err.Error()))
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	user, err := ah.authUseCases.Register.Execute(req.Username, req.Password)
	if err != nil {
		slog.Info("ошибка при регистрации пользователя", slog.String("error", err.Error()))
		http.Error(w, "Error on register user", http.StatusInternalServerError)
		return
	}

	slog.Info("зарегистрирован пользователь", slog.Any("User", user))
	w.WriteHeader(http.StatusCreated)
}

func (ah *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandlers) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
}
