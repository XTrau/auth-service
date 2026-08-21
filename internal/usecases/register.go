package usecases

import (
	"context"
	"errors"

	"github.com/XTrau/auth-service/internal/domain"
)

type RegisterUseCase struct {
	passwordHasher domain.Hasher
	unitOfWork     domain.UnitOfWork
}

func NewRegisterUseCase(hasher domain.Hasher, unitOfWork domain.UnitOfWork) *RegisterUseCase {
	return &RegisterUseCase{
		passwordHasher: hasher,
		unitOfWork:     unitOfWork,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, username string, password string) (user *domain.User, err error) {
	err = uc.unitOfWork.Execute(ctx, func(ctx context.Context, repos domain.Repositories) error {
		// Проверим существует ли пользователь с таким username
		u, err := repos.Users().GetByUsername(ctx, username)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return err
		}

		if err == nil && u != nil {
			return domain.ErrUsernameAlreadyExists
		}

		// Хешируем пароль
		passwordHash, err := uc.passwordHasher.Hash(password)

		if err != nil {
			return err
		}

		// Создаем пользователя
		user, err = repos.Users().Create(ctx, username, passwordHash)

		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}
