package repositories

import (
	"context"
	"database/sql"
	"time"
)

type PostgresBlockedTokensRepository struct {
	tx *sql.Tx
}

func NewPostgresBlockedTokensRepository(tx *sql.Tx) *PostgresBlockedTokensRepository {
	return &PostgresBlockedTokensRepository{tx}
}

func (repo *PostgresBlockedTokensRepository) Create(ctx context.Context, tokenString string, expiresAt time.Time) error {
	query := "INSERT INTO blocked_tokens (token, expires_at) VALUES ($1, $2)"

	_, err := repo.tx.ExecContext(ctx, query, tokenString, expiresAt)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresBlockedTokensRepository) Find(ctx context.Context, tokenString string) (bool, error) {
	query := "SELECT EXISTS(SELECT * FROM blocked_tokens WHERE token=$1)"

	var found bool

	if err := repo.tx.QueryRowContext(ctx, query, tokenString).Scan(&found); err != nil {
		return false, err
	}

	return found, nil
}
