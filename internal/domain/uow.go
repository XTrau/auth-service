package domain

import (
	"context"
)

type Repositories interface {
	Users() UserRepository
	BlockedTokens() BlockedTokensRepository
}

type UnitOfWork interface {
	// Выполняет функцию в рамках одной транзакции к бд
	Execute(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
	// Выполняет функцию в рамках одной транзакции к бд, при ошибках бд повторяет попытку
	ExecuteWithRetry(
		ctx context.Context,
		attempts int,
		fn func(ctx context.Context, repos Repositories) error,
	) error
}
