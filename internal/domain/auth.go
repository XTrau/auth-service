package domain

import (
	"time"
)

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenPayload struct {
	Subject   UserID `json:"id"`
	Username  string `json:"username"`
	ExpiresAt time.Time
}

type TokenPair struct {
	Access  string
	Refresh string
}

type TokenEncoder interface {
	Encode(payload TokenPayload, tokenType string, expire time.Duration) (string, error)
}

type TokenGenerator interface {
	Generate(payload TokenPayload) (TokenPair, error)
}

type TokenVerifier interface {
	Verify(token string) error
}

type TokenDecoder interface {
	Decode(token string, tokenType string) (TokenPayload, error)
}
