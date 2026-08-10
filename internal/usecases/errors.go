package usecases

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("User with that username already exists")
)
