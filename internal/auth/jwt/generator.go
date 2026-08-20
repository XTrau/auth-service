package jwt

import (
	"time"

	"github.com/XTrau/auth-service/internal/domain"
)

const (
	AccessType  = "access"
	RefreshType = "refresh"
)

type RS256Generator struct {
	encoder *RS256Encoder
}

func NewJwtGenerator(encoder *RS256Encoder) *RS256Generator {
	return &RS256Generator{
		encoder: encoder,
	}
}

func (g *RS256Generator) generateAccessToken(payload domain.TokenPayload) (string, error) {
	access, err := g.encoder.Encode(payload, AccessType, time.Minute*15)
	if err != nil {
		return "", err
	}
	return access, nil
}

func (g *RS256Generator) generateRefreshToken(payload domain.TokenPayload) (string, error) {
	refresh, err := g.encoder.Encode(payload, RefreshType, time.Hour*24*30)
	if err != nil {
		return "", err
	}
	return refresh, nil
}

func (g *RS256Generator) Generate(payload domain.TokenPayload) (domain.TokenPair, error) {
	access, err := g.generateAccessToken(payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refresh, err := g.generateRefreshToken(payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}
