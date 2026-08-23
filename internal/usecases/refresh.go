package usecases

import (
	"context"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type RefreshUseCase struct {
	generator  domain.TokenGenerator
	decoder    domain.TokenDecoder
	unitOfWork domain.UnitOfWork
}

func NewRefreshUseCase(generator domain.TokenGenerator, decoder domain.TokenDecoder, unitOfWork domain.UnitOfWork) *RefreshUseCase {
	return &RefreshUseCase{
		generator:  generator,
		decoder:    decoder,
		unitOfWork: unitOfWork,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, refreshToken string) (pair domain.TokenPair, err error) {
	const attempts int = 3

	err = uc.unitOfWork.ExecuteWithRetry(ctx, attempts, func(ctx context.Context, repos domain.Repositories) error {
		// Проверить токен в заблокированных
		found, err := repos.BlockedTokens().Find(ctx, refreshToken)
		if err != nil {
			return err
		}

		if found {
			return domain.ErrTokenBlocked
		}

		// Декодируем токен
		payload, err := uc.decoder.Decode(refreshToken, authjwt.RefreshType)
		if err != nil {
			return err
		}

		// Генерируем новую пару токенов
		pair, err = uc.generator.Generate(payload)
		if err != nil {
			return err
		}

		// Блокируем старый токен
		err = repos.BlockedTokens().Create(ctx, refreshToken, payload.ExpiresAt.UTC())
		if err != nil {
			return err
		}

		return nil
	})

	return pair, err

}
