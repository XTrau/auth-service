package usecases

import (
	"context"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type BlockTokenUseCase struct {
	decoder    domain.TokenDecoder
	unitOfWork domain.UnitOfWork
}

func NewBlockTokenUseCase(decoder domain.TokenDecoder, unitOfWork domain.UnitOfWork) *BlockTokenUseCase {
	return &BlockTokenUseCase{
		decoder:    decoder,
		unitOfWork: unitOfWork,
	}
}

func (uc *BlockTokenUseCase) Execute(ctx context.Context, refreshTokenString string) error {
	err := uc.unitOfWork.Execute(ctx, func(ctx context.Context, repos domain.Repositories) error {
		// декодировать токен (получить expires_at)
		payload, err := uc.decoder.Decode(refreshTokenString, authjwt.RefreshType)
		if err != nil {
			return err
		}

		// сохранить в заблокированные
		err = repos.BlockedTokens().Create(ctx, refreshTokenString, payload.ExpiresAt.UTC())
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
