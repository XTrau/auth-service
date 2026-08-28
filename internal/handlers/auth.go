package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/XTrau/auth-service/internal/domain"
	"github.com/XTrau/auth-service/internal/dto"
	errs "github.com/XTrau/auth-service/internal/errors"
	"github.com/asaskevich/govalidator"
)

const RefreshTokenName = "refresh"

type AuthorizationService interface {
	RegisterUser(ctx context.Context, username, password string) (user *domain.User, err error)
	LoginUser(ctx context.Context, username, password string) (pair domain.TokenPair, err error)
	BlockRefreshToken(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (pair domain.TokenPair, err error)
	GetUser(ctx context.Context, accessToken string) (user *domain.User, err error)
}

type AuthHandlers struct {
	service AuthorizationService
}

func NewAuthHandlers(service AuthorizationService) *AuthHandlers {
	return &AuthHandlers{service}
}

// RegisterHandler    godoc
// @Summary           Регистрация пользователя
// @Description       Регистрирует пользователя в системе
// @Description       логиниться в системе нужно отдельно от регистрации
// @Tags              Аутентификация
// @Accept            json
// @Param             userObject body dto.RegisterRequest true "userData"
// @Success           201
// @Router            /auth/register [post]
func (ah *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("плохое тело запроса на /register", slog.String("error", err.Error()))
		http.Error(w, errs.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if ok, err := govalidator.ValidateStruct(req); !ok || err != nil {
		slog.Info("ошибка при валидации тела /register", slog.Any("username", req.Username), slog.String("error", err.Error()))
		http.Error(w, errs.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	user, err := ah.service.RegisterUser(r.Context(), req.Username, req.Password)

	if errors.Is(err, errs.ErrUsernameAlreadyExists) {
		slog.Info("попытка зарегистрировать существующего пользователя", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		slog.Error("ошибка при регистрации пользователя", slog.String("error", err.Error()))
		http.Error(w, "Error on register user", http.StatusInternalServerError)
		return
	}

	slog.Info("зарегистрирован пользователь",
		slog.Int64("ID", int64(user.ID)),
		slog.String("username", user.Username),
	)
	w.WriteHeader(http.StatusCreated)
}

// LoginHandler    godoc
// @Summary        Логин пользователя
// @Description    Возвращает access token в body и refresh token в cookie
// @Description    Для доступа к api нужен access token
// @Description    Время жизни access token - 15 минут
// @Description    Время жизни refresh token - 30 дней
// @Description    Для обновления access token необходимо делать запрос POST /refresh с refresh token в куках
// @Tags           Аутентификация
// @Accept         json
// @Param          userObject body dto.LoginRequest true "userData"
// @Success        200 {object} handlers.AuthHandlers.LoginHandler.response
// @Router         /auth/login [post]
func (ah *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("плохое тело запроса на /login", slog.String("error", err.Error()))
		http.Error(w, errs.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := ah.service.LoginUser(r.Context(), req.Login, req.Password)
	if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrInvalidPassword) {
		http.Error(w, "incorrect username or password", http.StatusBadRequest)
		return
	}

	if err != nil {
		slog.Error("ошибка при логине пользователя", slog.String("error", err.Error()))
		http.Error(w, "error on user login", http.StatusInternalServerError)
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
		Path:     "/auth",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refreshCookie)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// LogoutHandler   godoc
// @Summary        Разлогинить пользователя
// @Description    Удаляет refresh token из cookie
// @Tags           Аутентификация
// @Accept         json
// @Success        204
// @Router         /auth/logout [post]
func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(RefreshTokenName)
	if err == http.ErrNoCookie {
		return
	}

	ah.service.BlockRefreshToken(r.Context(), c.Value)

	cookie := &http.Cookie{
		Name:     RefreshTokenName,
		Value:    "",
		Path:     "/auth",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

// RefreshTokensHandler godoc
// @Summary        		Обновить токены
// @Description    		Возвращает новый access и refresh токены
// @Description    		Требуется refresh токен в куках
// @Tags           		Аутентификация
// @Success        		200
// @Router         		/auth/refresh [post]
func (ah *AuthHandlers) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(RefreshTokenName)
	if err == http.ErrNoCookie {
		http.Error(w, "user not logged in", http.StatusUnauthorized)
		return
	}

	// Создаем новую пару токенов и блокируем старый
	tokens, err := ah.service.RefreshTokens(r.Context(), c.Value)
	if err != nil {
		slog.Info("плохой jwt токен пользователя", slog.String("error", err.Error()))
		http.Error(w, "user not logged in", http.StatusUnauthorized)
		return
	}

	b, err := json.Marshal(dto.AccessTokenResponse{Token: tokens.Access})
	if err != nil {
		slog.Info("ошибка при marshal access token", slog.String("error", err.Error()))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     RefreshTokenName,
		Value:    tokens.Refresh,
		Path:     "/auth",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refreshCookie)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// GetUserHandler godoc
// @Summary       Получить пользователя
// @Description   Возвращает данные пользователя по access_token
// @Tags          Аутентификация
// @Param 		  Authorization header string true "Example: Bearer {token}"
// @Success       200   {object} dto.UserDataResponse
// @Success	 	  401	{string} string "unauthorized user"
// @Router        /auth/user [get]
func (ah *AuthHandlers) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := ah.service.GetUser(r.Context(), accessTokenString)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userResp := dto.UserDataResponse{
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		slog.Info("ошибка при marshal userResp", slog.String("error", err.Error()))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}
