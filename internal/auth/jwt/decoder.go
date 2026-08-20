package jwt

import (
	"crypto/rsa"
	"strconv"

	"github.com/XTrau/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type RS256Decoder struct {
	publicKey *rsa.PublicKey
}

func NewRS256Decoder(cfg RS256Config) *RS256Decoder {
	return &RS256Decoder{
		publicKey: cfg.PublicRSAKey(),
	}
}

func (d *RS256Decoder) parseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return d.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, ErrInvalidToken
	}

	return t, nil
}

func (d *RS256Decoder) Verify(tokenString string) error {
	t, err := d.parseToken(tokenString, nil)

	if err != nil {
		return err
	}

	if t.Method != jwt.SigningMethodRS256 {
		return ErrInvalidSigningMethod
	}

	return nil
}

func (d *RS256Decoder) Decode(tokenString string, tokenType string) (domain.TokenPayload, error) {
	claims := new(CustomClaims)
	t, err := d.parseToken(tokenString, claims)

	if err != nil {
		return domain.TokenPayload{}, err
	}

	if t.Method != jwt.SigningMethodRS256 {
		return domain.TokenPayload{}, ErrInvalidSigningMethod
	}

	if claims.Type != tokenType {
		return domain.TokenPayload{}, ErrTokenType
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 0)
	if err != nil {
		return domain.TokenPayload{}, err
	}

	result := domain.TokenPayload{
		Subject:   domain.UserID(userID),
		Username:  claims.Username,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return result, nil
}
