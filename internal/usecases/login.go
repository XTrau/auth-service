package usecases

import (
	"context"

	"github.com/XTrau/auth-service/internal/domain"
)

type LoginUseCase struct {
	tokenGenerator domain.TokenGenerator
	passwordHasher domain.Hasher
	userRepository domain.UserRepository
}

func NewLoginUseCase(generator domain.TokenGenerator, hasher domain.Hasher, userRepository domain.UserRepository) *LoginUseCase {
	return &LoginUseCase{
		tokenGenerator: generator,
		passwordHasher: hasher,
		userRepository: userRepository,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, username, password string) (domain.TokenPair, error) {
	user, err := uc.userRepository.GetByUsername(ctx, username)

	if err != nil {
		return domain.TokenPair{}, err
	}

	if user == nil {
		return domain.TokenPair{}, ErrUserNotFound
	}

	if !uc.passwordHasher.Compare(user.PasswordHash, password) {
		return domain.TokenPair{}, ErrInvalidPassword
	}

	payload := domain.TokenPayload{
		Subject:  user.ID,
		Username: user.Username,
	}

	return uc.tokenGenerator.Generate(payload)
}
