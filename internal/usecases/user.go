package usecases

import (
	"context"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type GetUserUseCase struct {
	tokenDecoder domain.TokenDecoder
	unitOfWork   domain.UnitOfWork
}

func NewGetUserUseCase(decoder domain.TokenDecoder, unitOfWork domain.UnitOfWork) *GetUserUseCase {
	return &GetUserUseCase{
		tokenDecoder: decoder,
		unitOfWork:   unitOfWork,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, accessToken string) (user *domain.User, err error) {
	err = uc.unitOfWork.Execute(ctx, func(ctx context.Context, repos domain.Repositories) error {
		// Достаем user id из токена
		payload, err := uc.tokenDecoder.Decode(accessToken, authjwt.AccessType)
		if err != nil {
			return err
		}

		// Находим пользователя по id
		user, err = repos.Users().GetByID(ctx, payload.Subject)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}
