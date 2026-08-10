package usecases

import "github.com/XTrau/auth-service/internal/domain"

type RegisterUseCase struct {
	passwordHasher domain.Hasher
	userRepository domain.UserRepository
}

func NewRegisterUseCase(hasher domain.Hasher, userRepository domain.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		passwordHasher: hasher,
		userRepository: userRepository,
	}
}

func (uc *RegisterUseCase) Execute(username string, password string) (*domain.User, error) {
	user, err := uc.userRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, ErrUsernameAlreadyExists
	}

	passwordHash, err := uc.passwordHasher.Hash(password)

	if err != nil {
		return nil, err
	}

	user, err = uc.userRepository.Create(username, passwordHash)

	if err != nil {
		return nil, err
	}

	return user, nil
}
