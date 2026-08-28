package repositories

import (
	"context"
	"time"
)

type PgBlockedTokensRepository struct {
}

func NewPgBlockedTokensRepository() *PgBlockedTokensRepository {
	return &PgBlockedTokensRepository{}
}

func (repo *PgBlockedTokensRepository) Create(ctx context.Context, tokenString string, expiresAt time.Time) error {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return err
	}

	query := "INSERT INTO blocked_tokens (token, expires_at) VALUES ($1, $2)"

	_, err = tx.ExecContext(ctx, query, tokenString, expiresAt)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PgBlockedTokensRepository) Find(ctx context.Context, tokenString string) (bool, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := "SELECT EXISTS(SELECT * FROM blocked_tokens WHERE token=$1)"

	var found bool

	if err := tx.QueryRowContext(ctx, query, tokenString).Scan(&found); err != nil {
		return false, err
	}

	return found, nil
}
