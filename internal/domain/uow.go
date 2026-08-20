package domain

import "context"

type Repositories interface {
	Users() UserRepository
	BlockedTokens() BlockedTokensRepository
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}
