package jwt

import (
	"time"

	"github.com/XTrau/auth-service/internal/domain"
)

const (
	AccessType  = "access"
	RefreshType = "refresh"
)

type JwtGenerator struct {
	encoder *JwtEncoder
}

func NewJwtGenerator(encoder *JwtEncoder) *JwtGenerator {
	return &JwtGenerator{
		encoder: encoder,
	}
}

func (jg *JwtGenerator) generateAccessToken(payload domain.TokenPayload) (string, error) {
	access, err := jg.encoder.Encode(payload, AccessType, time.Minute*15)
	if err != nil {
		return "", err
	}
	return access, nil
}

func (jg *JwtGenerator) generateRefreshToken(payload domain.TokenPayload) (string, error) {
	refresh, err := jg.encoder.Encode(payload, RefreshType, time.Hour*24*30)
	if err != nil {
		return "", err
	}
	return refresh, nil
}

func (jg *JwtGenerator) Generate(payload domain.TokenPayload) (domain.TokenPair, error) {
	access, err := jg.generateAccessToken(payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refresh, err := jg.generateRefreshToken(payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}
