package service_test

import (
	"context"
	"errors"
	"testing"

	errs "github.com/XTrau/auth-service/internal/errors"
	"github.com/XTrau/auth-service/internal/service"
)

func TestRegisterUserValidation(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		password      string
		expectedField string
		expectedError error
	}{
		{
			"Регистрация нормального пользователя",
			"username",
			"password",
			"",
			nil,
		},
		{
			"Регистрация пользователя с коротким ником",
			"user",
			"password",
			"username",
			errs.ErrUsernameTooShort,
		},
		{
			"Регистрация пользователя с слишком длинным ником (65 символов)",
			"usernameusernameusernameusernameusernameusernameusernameusername1",
			"password",
			"username",
			errs.ErrUsernameTooLong,
		},
		{
			"Регистрация пользователя с коротким паролем",
			"username",
			"pass",
			"password",
			errs.ErrPasswordTooShort,
		},
		{
			"Регистрация пользователя с слишком длинным паролем (65 символов)",
			"username",
			"usernameusernameusernameusernameusernameusernameusernameusername1",
			"password",
			errs.ErrPasswordTooLong,
		},
	}

	tokenizer := &MockTokenizer{}
	hasher := &MockHasher{}
	unitOfWork := &MockUnitOfWork{}
	userRepo := NewMockUserRepository()
	blockedTokensRepo := NewMockBlockedTokensRepository()

	service := service.NewAuthenticationService(tokenizer, hasher, unitOfWork, userRepo, blockedTokensRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.RegisterUser(context.Background(), tc.username, tc.password)

			if tc.expectedError == nil {
				if err != nil {
					t.Fatalf("Неожиданная ошибка: %v, ожидаемая ошибка: %v", err, tc.expectedError)
				}
				return
			}

			var validationErr errs.ValidationError

			if !errors.As(err, &validationErr) {
				t.Fatalf("Ожидалась ValidationError получена: %v", err)
				return
			}

			if validationErr.Field != tc.expectedField {
				t.Fatalf("Ошибка по неожиданному полю: %v, ожидаемое поле: %v", validationErr.Field, tc.expectedField)
				return
			}

			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("Неожиданная ошибка: %v, ожидаемая ошибка: %v", err, tc.expectedError)
				return
			}
		})
	}
}
