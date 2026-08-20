package usecases

import (
	"context"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type BlockTokenUseCase struct {
	blockedTokensRepository domain.BlockedTokensRepository
	decoder                 domain.TokenDecoder
}

func NewBlockTokenUseCase(blockedTokensRepository domain.BlockedTokensRepository, decoder domain.TokenDecoder) *BlockTokenUseCase {
	return &BlockTokenUseCase{
		blockedTokensRepository: blockedTokensRepository,
		decoder:                 decoder,
	}
}

func (uc *BlockTokenUseCase) Execute(ctx context.Context, refreshTokenString string) error {
	// декодировать токен (получить expires_at)
	payload, err := uc.decoder.Decode(refreshTokenString, authjwt.RefreshType)
	if err != nil {
		return err
	}

	// сохранить в заблокированные
	err = uc.blockedTokensRepository.Create(ctx, refreshTokenString, payload.ExpiresAt.UTC())
	if err != nil {
		return err
	}

	return nil
}
