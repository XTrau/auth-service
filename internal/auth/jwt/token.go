package jwt

import (
	"crypto/rsa"
	"strconv"
	"time"

	"github.com/XTrau/auth-service/internal/domain"
	errs "github.com/XTrau/auth-service/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessType  = "access"
	RefreshType = "refresh"
)

type RS256Config interface {
	PrivateRSAKey() *rsa.PrivateKey
	PublicRSAKey() *rsa.PublicKey
}

type RS256Tokenizer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRS256Tokenizer(cfg RS256Config) *RS256Tokenizer {
	return &RS256Tokenizer{
		privateKey: cfg.PrivateRSAKey(),
		publicKey:  cfg.PublicRSAKey(),
	}
}

type CustomClaims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// Кодирует payload в jwt token
func (t *RS256Tokenizer) encode(payload domain.TokenPayload, tokenType string, expire time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, CustomClaims{
		Username: payload.Username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.Subject.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Генерирует access jwt токен для payload
func (t *RS256Tokenizer) generateAccessToken(payload domain.TokenPayload) (string, error) {
	access, err := t.encode(payload, AccessType, time.Minute*15)
	if err != nil {
		return "", err
	}
	return access, nil
}

// Генерирует refresh jwt токен для payload
func (t *RS256Tokenizer) generateRefreshToken(payload domain.TokenPayload) (string, error) {
	refresh, err := t.encode(payload, RefreshType, time.Hour*24*30)
	if err != nil {
		return "", err
	}
	return refresh, nil
}

// Генерирует пару jwt токенов для payload
func (t *RS256Tokenizer) Generate(payload domain.TokenPayload) (domain.TokenPair, error) {
	access, err := t.generateAccessToken(payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refresh, err := t.generateRefreshToken(payload)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}

// Парсит токен, возвращает ошибку если токен не валидный
func (t *RS256Tokenizer) parseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return t.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return token, nil
}

// Достает TokenPayload из токена, возвращает ошибку если токен не валидный
func (t *RS256Tokenizer) Decode(tokenString string, tokenType string) (domain.TokenPayload, error) {
	claims := new(CustomClaims)
	token, err := t.parseToken(tokenString, claims)

	if err != nil {
		return domain.TokenPayload{}, err
	}

	if token.Method != jwt.SigningMethodRS256 {
		return domain.TokenPayload{}, errs.ErrInvalidSigningMethod
	}

	if claims.Type != tokenType {
		return domain.TokenPayload{}, errs.ErrTokenType
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
