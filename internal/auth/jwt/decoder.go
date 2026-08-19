package jwt

import (
	"crypto/rsa"
	"strconv"

	"github.com/XTrau/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JwtDecoder struct {
	publicKey *rsa.PublicKey
}

func NewJwtDecoder(cfg RSA256Config) *JwtDecoder {
	return &JwtDecoder{
		publicKey: cfg.PublicRSAKey(),
	}
}

func (jd *JwtDecoder) parseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jd.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, ErrInvalidToken
	}

	return t, nil
}

func (jd *JwtDecoder) Verify(tokenString string) error {
	_, err := jd.parseToken(tokenString, nil)
	return err
}

func (jd *JwtDecoder) Decode(tokenString string, tokenType string) (domain.TokenPayload, error) {
	claims := new(CustomClaims)
	_, err := jd.parseToken(tokenString, claims)

	if err != nil {
		return domain.TokenPayload{}, err
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
