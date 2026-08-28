package errors

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("Пользователь с таким username уже существует")
	ErrUserNotFound          = errors.New("Пользователь не найден")
	ErrInvalidPassword       = errors.New("Неправильный пароль")
	ErrTokenBlocked          = errors.New("Токен заблокирован")
)

var (
	ErrUsernameTooShort            = NewValidationError("username", "Username слишком короткий, минимальный размер - 8")
	ErrUsernameTooLong             = NewValidationError("username", "Username слишком длинный, максимальный размер - 64")
	ErrUsenameBadFirstSymbol       = NewValidationError("username", "Username может начинаться только с буквы или нижнего подчеркивания")
	ErrUsenameUnresolvedCharacters = NewValidationError("username", "Username может содержать только символы английского алфавита, числа и символ нижнего подчеркивания")
	ErrPasswordTooShort            = NewValidationError("password", "Пароль слишком короткий, минимальный размер - 8")
	ErrPasswordTooLong             = NewValidationError("password", "Пароль слишком длинный, максимальный размер - 64")
)
