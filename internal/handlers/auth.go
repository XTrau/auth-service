package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
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

// LogoutHandler   godoc
// @Summary        Разлогинить пользователя
// @Description    Удаляет refresh token из cookie
// @Tags           Аутентификация
// @Accept         json
// @Success        204
// @Router         /auth/logout [post]
func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(RefreshTokenName)
	if err == http.ErrNoCookie {
		http.Error(w, "user not logged in", http.StatusBadRequest)
		return
	}

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

	tokens, err := ah.authUseCases.Refresh.Execute(c.Value)
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
// @Param 		  Authorization header string true "Example: Bearer "
// @Success       200   {object} dto.UserDataResponse
// @Success	 	  401	{string} string "unauthorized user"
// @Router        /auth/user [get]
func (ah *AuthHandlers) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") && len(authHeader) > 7 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessTokenString := authHeader[7:]
	user, err := ah.authUseCases.User.Execute(accessTokenString)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userResp := dto.UserDataResponse{
		Username: user.Username,
	}

	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		slog.Info("ошибка при marshal userResp", slog.String("error", err.Error()))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
}
