package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/XTrau/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type RS256Encoder struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRS256Encoder(cfg RS256Config) *RS256Encoder {
	return &RS256Encoder{
		privateKey: cfg.PrivateRSAKey(),
		publicKey:  cfg.PublicRSAKey(),
	}
}

type CustomClaims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

func (e *RS256Encoder) Encode(payload domain.TokenPayload, tokenType string, expire time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, CustomClaims{
		Username: payload.Username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.Subject.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(e.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
