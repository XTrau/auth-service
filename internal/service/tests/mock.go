package service_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/XTrau/auth-service/internal/domain"
	errs "github.com/XTrau/auth-service/internal/errors"
)

//	type Tokenizer interface {
//		Generate(payload domain.TokenPayload) (domain.TokenPair, error)
//		Decode(token, tokenType string) (domain.TokenPayload, error)
//	}
type MockTokenizer struct {
}

func (t *MockTokenizer) Generate(payload domain.TokenPayload) (domain.TokenPair, error) {
	token, _ := json.Marshal(payload)

	return domain.TokenPair{
		Access:  string(token),
		Refresh: string(token),
	}, nil
}

func (t *MockTokenizer) Decode(token, tokenType string) (domain.TokenPayload, error) {
	var payload domain.TokenPayload

	json.Unmarshal([]byte(token), &payload)
	payload.ExpiresAt = time.Now().Add(time.Hour)

	return payload, nil
}

//	type Hasher interface {
//		Hash(password string) (string, error)
//		Compare(hash, password string) bool
//	}
type MockHasher struct {
}

func (h *MockHasher) Hash(password string) (string, error) {
	return password, nil
}

func (h *MockHasher) Compare(hash, password string) bool {
	return hash == password
}

// type BlockedTokensRepository interface {
// 	Create(ctx context.Context, tokenString string, expiresAt time.Time) error
// 	Find(ctx context.Context, tokenString string) (bool, error)
// }

type MockBlockedTokensRepository struct {
	blockedTokens map[string]struct{}
}

func NewMockBlockedTokensRepository() *MockBlockedTokensRepository {
	return &MockBlockedTokensRepository{
		blockedTokens: make(map[string]struct{}),
	}
}

func (repo *MockBlockedTokensRepository) Create(ctx context.Context, tokenString string, expiresAt time.Time) error {
	repo.blockedTokens[tokenString] = struct{}{}
	return nil
}

func (repo *MockBlockedTokensRepository) Find(ctx context.Context, tokenString string) (bool, error) {
	_, found := repo.blockedTokens[tokenString]
	return found, nil
}

// type UserRepository interface {
// 	Create(ctx context.Context, user domain.User) (domain.User, error)
// 	GetByID(ctx context.Context, id int64) (domain.User, error)
// 	GetByUsername(ctx context.Context, username string) (domain.User, error)
// }

type MockUserRepository struct {
	users  map[int64]domain.User
	nextID int64
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[int64]domain.User),
		nextID: 1,
	}
}

func (repo *MockUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = repo.nextID
	repo.users[repo.nextID] = user
	repo.nextID++
	return user, nil
}

func (repo *MockUserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	user, found := repo.users[repo.nextID]
	if !found {
		return domain.User{}, errs.ErrUserNotFound
	}
	return user, nil
}

func (repo *MockUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	for _, user := range repo.users {
		if user.Username == username {
			return user, nil
		}
	}
	return domain.User{}, errs.ErrUserNotFound
}

// type UnitOfWork interface {
// 	// Выполняет функцию в рамках одной транзакции, при ошибках связанных с хранилищем повторяет попытку
// 	ExecuteWithRetry(ctx context.Context, attempts int, fn func(ctx context.Context) error) error
// }

type MockUnitOfWork struct {
}

func (uow *MockUnitOfWork) ExecuteWithRetry(ctx context.Context, attempts int, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
