package domain

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("User with that username already exists")
	ErrUserNotFound          = errors.New("User not found")
	ErrInvalidPassword       = errors.New("Invalid password")
	ErrTokenBlocked          = errors.New("Token blocked")
)
