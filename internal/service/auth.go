package service

import (
	"context"
	"errors"
	"time"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
	errs "github.com/XTrau/auth-service/internal/errors"
)

var defaultAttempts = 3

type Tokenizer interface {
	Generate(payload domain.TokenPayload) (domain.TokenPair, error)
	Decode(token, tokenType string) (domain.TokenPayload, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type BlockedTokensRepository interface {
	Create(ctx context.Context, tokenString string, expiresAt time.Time) error
	Find(ctx context.Context, tokenString string) (bool, error)
}

type UserRepository interface {
	Create(ctx context.Context, username string, passwordHash string) (*domain.User, error)
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type UnitOfWork interface {
	// Выполняет функцию в рамках одной транзакции, при ошибках связанных с хранилищем повторяет попытку
	ExecuteWithRetry(ctx context.Context, attempts int, fn func(ctx context.Context) error) error
}

type authorizationService struct {
	hasher            Hasher
	tokenizer         Tokenizer
	unitOfWork        UnitOfWork
	userRepo          UserRepository
	blockedTokensRepo BlockedTokensRepository
}

func NewAuthorizationService(
	tokenizer Tokenizer,
	hasher Hasher,
	unitOfWork UnitOfWork,
	userRepo UserRepository,
	blockedTokensRepo BlockedTokensRepository,
) *authorizationService {
	return &authorizationService{
		hasher:            hasher,
		tokenizer:         tokenizer,
		unitOfWork:        unitOfWork,
		userRepo:          userRepo,
		blockedTokensRepo: blockedTokensRepo,
	}
}

func (s *authorizationService) RegisterUser(ctx context.Context, username, password string) (user *domain.User, err error) {
	err = s.unitOfWork.ExecuteWithRetry(ctx, defaultAttempts, func(ctx context.Context) error {
		// Проверим существует ли пользователь с таким username
		u, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
			return err
		}

		if err == nil && u != nil {
			return errs.ErrUsernameAlreadyExists
		}

		// Хешируем пароль
		passwordHash, err := s.hasher.Hash(password)

		if err != nil {
			return err
		}

		// Создаем пользователя
		user, err = s.userRepo.Create(ctx, username, passwordHash)

		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func (s *authorizationService) LoginUser(ctx context.Context, username, password string) (pair domain.TokenPair, err error) {
	err = s.unitOfWork.ExecuteWithRetry(ctx, defaultAttempts, func(ctx context.Context) error {
		// Находим пользователя
		user, err := s.userRepo.GetByUsername(ctx, username)

		if err != nil {
			return err
		}

		// Проверям пароль с хешем
		if !s.hasher.Compare(user.PasswordHash, password) {
			return errs.ErrInvalidPassword
		}

		payload := domain.TokenPayload{
			Subject:  user.ID,
			Username: user.Username,
		}

		// Генерируем пару токенов
		pair, err = s.tokenizer.Generate(payload)
		if err != nil {
			return err
		}

		return nil
	})

	return pair, err
}

func (s *authorizationService) BlockRefreshToken(ctx context.Context, refreshToken string) error {
	err := s.unitOfWork.ExecuteWithRetry(ctx, defaultAttempts, func(ctx context.Context) error {
		// декодировать токен (получить expires_at)
		payload, err := s.tokenizer.Decode(refreshToken, authjwt.RefreshType)
		if err != nil {
			return err
		}

		// сохранить в заблокированные
		err = s.blockedTokensRepo.Create(ctx, refreshToken, payload.ExpiresAt.UTC())
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *authorizationService) RefreshTokens(ctx context.Context, refreshToken string) (pair domain.TokenPair, err error) {
	err = s.unitOfWork.ExecuteWithRetry(ctx, defaultAttempts, func(ctx context.Context) error {
		// Проверить токен в заблокированных
		found, err := s.blockedTokensRepo.Find(ctx, refreshToken)
		if err != nil {
			return err
		}

		if found {
			return errs.ErrTokenBlocked
		}

		// Декодируем токен
		payload, err := s.tokenizer.Decode(refreshToken, authjwt.RefreshType)
		if err != nil {
			return err
		}

		// Генерируем новую пару токенов
		pair, err = s.tokenizer.Generate(payload)
		if err != nil {
			return err
		}

		// Блокируем старый токен
		err = s.blockedTokensRepo.Create(ctx, refreshToken, payload.ExpiresAt.UTC())
		if err != nil {
			return err
		}

		return nil
	})

	return pair, err
}

func (s *authorizationService) GetUser(ctx context.Context, accessToken string) (user *domain.User, err error) {
	// Декодируем токен, получаем id
	payload, err := s.tokenizer.Decode(accessToken, authjwt.AccessType)
	if err != nil {
		return nil, err
	}

	err = s.unitOfWork.ExecuteWithRetry(ctx, defaultAttempts, func(ctx context.Context) error {
		// Находим пользователя по id
		user, err = s.userRepo.GetByID(ctx, payload.Subject)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}
