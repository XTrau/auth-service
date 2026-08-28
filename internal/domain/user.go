package domain

import (
	"unicode"
	"unicode/utf8"

	errs "github.com/XTrau/auth-service/internal/errors"
)

type username string

func isAllowedSymbol(r rune) bool {
	return (unicode.IsLetter(r) && r < unicode.MaxASCII) || unicode.IsDigit(r) || r == '_'
}

func NewUsername(value string) (username, error) {
	length := utf8.RuneCountInString(value)

	if length < 8 {
		return "", errs.ErrUsernameTooShort
	}

	if length > 64 {
		return "", errs.ErrUsernameTooLong
	}

	if !unicode.IsLetter(rune(value[0])) && rune(value[0]) != '_' {
		return "", errs.ErrUsenameBadFirstSymbol
	}

	for _, r := range value {
		if !isAllowedSymbol(r) {
			return "", errs.ErrUsenameUnresolvedCharacters
		}
	}

	return username(value), nil
}

type password string

func NewPassword(value string) (password, error) {
	length := utf8.RuneCountInString(value)
	if length < 8 {
		return "", errs.ErrPasswordTooShort
	}

	if length > 64 {
		return "", errs.ErrPasswordTooLong
	}

	return password(value), nil
}

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}
