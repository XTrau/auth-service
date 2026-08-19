package repositories

import (
	"context"
	"database/sql"
	"time"
)

type PostgresBlockedTokensRepository struct {
	db *sql.DB
}

func NewPostgresBlockedTokensRepository(db *sql.DB) *PostgresBlockedTokensRepository {
	return &PostgresBlockedTokensRepository{db}
}

func (repo *PostgresBlockedTokensRepository) Create(ctx context.Context, tokenString string, expiresAt time.Time) error {
	query := "INSERT INTO blocked_tokens (token, expires_at) VALUES ($1, $2)"

	_, err := repo.db.ExecContext(ctx, query, tokenString, expiresAt)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresBlockedTokensRepository) Find(ctx context.Context, tokenString string) (bool, error) {
	query := "SELECT EXISTS(SELECT * FROM blocked_tokens WHERE token=$1)"

	var found bool

	if err := repo.db.QueryRowContext(ctx, query, tokenString).Scan(&found); err != nil {
		return false, err
	}

	return found, nil
}
