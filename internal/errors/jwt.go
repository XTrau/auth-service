package errors

import "errors"

var (
	ErrInvalidToken         = errors.New("Invalid token")
	ErrBadJwtClaims         = errors.New("Bad jwt claims")
	ErrTokenType            = errors.New("Incorrect token type")
	ErrInvalidSigningMethod = errors.New("Invalid Signing Method")
)
