package usecases

import (
	"context"

	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type GetUserUseCase struct {
	tokenDecoder   domain.TokenDecoder
	userRepository domain.UserRepository
}

func NewGetUserUseCase(decoder domain.TokenDecoder, userRepository domain.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		tokenDecoder:   decoder,
		userRepository: userRepository,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, accessToken string) (*domain.User, error) {
	payload, err := uc.tokenDecoder.Decode(accessToken, authjwt.AccessType)
	if err != nil {
		return nil, err
	}

	return uc.userRepository.GetByID(ctx, payload.Subject)
}
