package domain

import (
	"context"
	"time"
)

type BlockedTokensRepository interface {
	Create(ctx context.Context, tokenString string, expiresAt time.Time) error
	Find(ctx context.Context, tokenString string) (bool, error)
}
