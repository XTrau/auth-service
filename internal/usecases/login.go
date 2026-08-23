package usecases

import (
	"context"

	"github.com/XTrau/auth-service/internal/domain"
)

type LoginUseCase struct {
	tokenGenerator domain.TokenGenerator
	passwordHasher domain.Hasher
	unitOfWork     domain.UnitOfWork
}

func NewLoginUseCase(generator domain.TokenGenerator, hasher domain.Hasher, unitOfWork domain.UnitOfWork) *LoginUseCase {
	return &LoginUseCase{
		tokenGenerator: generator,
		passwordHasher: hasher,
		unitOfWork:     unitOfWork,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, username, password string) (pair domain.TokenPair, err error) {
	const attempts int = 3
	
	err = uc.unitOfWork.ExecuteWithRetry(ctx, attempts, func(ctx context.Context, repos domain.Repositories) error {
		// Находим пользователя
		user, err := repos.Users().GetByUsername(ctx, username)

		if err != nil {
			return err
		}

		// Проверям пароль с хешем
		if !uc.passwordHasher.Compare(user.PasswordHash, password) {
			return domain.ErrInvalidPassword
		}

		payload := domain.TokenPayload{
			Subject:  user.ID,
			Username: user.Username,
		}

		// Генерируем пару токенов
		pair, err = uc.tokenGenerator.Generate(payload)
		if err != nil {
			return err
		}

		return nil
	})

	return pair, err
}
