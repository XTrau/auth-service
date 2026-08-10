package usecases

import (
	authjwt "github.com/XTrau/auth-service/internal/auth/jwt"
	"github.com/XTrau/auth-service/internal/domain"
)

type RefreshUseCase struct {
	generator domain.TokenGenerator
	decoder   domain.TokenDecoder
}

func NewRefreshUseCase(generator domain.TokenGenerator, decoder domain.TokenDecoder) *RefreshUseCase {
	return &RefreshUseCase{
		generator: generator,
		decoder:   decoder,
	}
}

func (uc *RefreshUseCase) Execute(refreshToken string) (domain.TokenPair, error) {
	payload, err := uc.decoder.Decode(refreshToken, authjwt.RefreshType)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return uc.generator.Generate(payload)
}
