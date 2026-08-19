package usecases

import (
	"context"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type LogoutUseCase struct {
	blockedTokensRepository domain.BlockedTokensRepository
	decoder                 domain.TokenDecoder
}

func NewLogoutUseCase(blockedTokensRepository domain.BlockedTokensRepository, decoder domain.TokenDecoder) *LogoutUseCase {
	return &LogoutUseCase{
		blockedTokensRepository: blockedTokensRepository,
		decoder:                 decoder,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, refreshTokenString string) error {
	// декодировать токен (получить expires_at)
	payload, err := uc.decoder.Decode(refreshTokenString, authjwt.RefreshType)
	if err != nil {
		return err
	}

	// сохранить в заблокированные
	err = uc.blockedTokensRepository.Create(ctx, refreshTokenString, payload.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}
