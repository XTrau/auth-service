package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/XTrau/auth-service/internal/dto"
	"github.com/XTrau/auth-service/internal/usecases"
	"github.com/asaskevich/govalidator"
)

const RefreshTokenName = "refresh"

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
	_, err := r.Cookie(RefreshTokenName)
	if err != http.ErrNoCookie {
		http.Error(w, "user already logged in", http.StatusBadRequest)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("плохое тело запроса на /login", slog.String("error", err.Error()))
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := ah.authUseCases.Login.Execute(req.Login, req.Password)
	if err != nil {
		slog.Info("ошибка при логине пользователя", slog.String("error", err.Error()))
		http.Error(w, "error on user login", http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(dto.AccessTokenResponse{Token: tokens.Access})
	if err != nil {
		slog.Info("ошибка при marshal access token", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     RefreshTokenName,
		Value:    tokens.Refresh,
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refreshCookie)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(RefreshTokenName)
	if err == http.ErrNoCookie {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     RefreshTokenName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func (ah *AuthHandlers) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(RefreshTokenName)
	if err == http.ErrNoCookie {
		http.Error(w, "user not logged in", http.StatusUnauthorized)
		return
	}

	tokens, err := ah.authUseCases.Refresh.Execute(c.Value)
	if err != nil {
		slog.Info("плохой jwt токен пользователя", slog.String("error", err.Error()))
		http.Error(w, "user not logged in", http.StatusUnauthorized)
		return
	}

	b, err := json.Marshal(dto.AccessTokenResponse{Token: tokens.Access})
	if err != nil {
		slog.Info("error on marshal access token", slog.String("error", err.Error()))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     RefreshTokenName,
		Value:    tokens.Refresh,
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refreshCookie)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (ah *AuthHandlers) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	
}
