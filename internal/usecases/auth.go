package usecases

import "github.com/XTrau/auth-service/internal/domain"

type AuthUseCases struct {
	Register *RegisterUseCase
	Login    *LoginUseCase
	Refresh  *RefreshUseCase
	User     *GetUserUseCase
}

func NewAuthUseCases(
	generator domain.TokenGenerator,
	decoder domain.TokenDecoder,
	hasher domain.Hasher,
	userRepository domain.UserRepository,
) *AuthUseCases {
	return &AuthUseCases{
		Register: NewRegisterUseCase(hasher, userRepository),
		Login:    NewLoginUseCase(generator, hasher, userRepository),
		Refresh:  NewRefreshUseCase(generator, decoder),
		User:     NewGetUserUseCase(decoder, userRepository),
	}
}
