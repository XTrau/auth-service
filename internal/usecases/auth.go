package usecases

import (
	"github.com/XTrau/auth-service/internal/domain"
)

type AuthUseCases struct {
	Register   *RegisterUseCase
	Login      *LoginUseCase
	TokenBlock *BlockTokenUseCase
	Refresh    *RefreshUseCase
	User       *GetUserUseCase
}

func NewAuthUseCases(
	generator domain.TokenGenerator,
	decoder domain.TokenDecoder,
	hasher domain.Hasher,
	unitOfWork domain.UnitOfWork,
) *AuthUseCases {
	return &AuthUseCases{
		Register:   NewRegisterUseCase(hasher, unitOfWork),
		Login:      NewLoginUseCase(generator, hasher, unitOfWork),
		TokenBlock: NewBlockTokenUseCase(decoder, unitOfWork),
		Refresh:    NewRefreshUseCase(generator, decoder, unitOfWork),
		User:       NewGetUserUseCase(decoder, unitOfWork),
	}
}
