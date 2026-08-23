package uow

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/XTrau/auth-service/internal/domain"
	"github.com/XTrau/auth-service/internal/repositories"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepositories struct {
	userRepository          *repositories.PostgresUserRepository
	blockedTokensRepository *repositories.PostgresBlockedTokensRepository
}

func (repos PostgresRepositories) Users() domain.UserRepository {
	return repos.userRepository
}

func (repos PostgresRepositories) BlockedTokens() domain.BlockedTokensRepository {
	return repos.blockedTokensRepository
}

type PostgresUnitOfWork struct {
	db *sql.DB
}

func NewPostgresUnitOfWork(db *sql.DB) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{db}
}

func (uow *PostgresUnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context, repos domain.Repositories) error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	repos := PostgresRepositories{
		userRepository:          repositories.NewPostgresUserRepository(tx),
		blockedTokensRepository: repositories.NewPostgresBlockedTokensRepository(tx),
	}

	if err := fn(ctx, repos); err != nil {
		return err
	}

	return tx.Commit()
}

func (uow *PostgresUnitOfWork) ExecuteWithRetry(
	ctx context.Context,
	attempts int,
	fn func(ctx context.Context, repos domain.Repositories) error,
) error {
	for i := 0; i < attempts; i++ {
		err := uow.Execute(ctx, fn)

		if err == nil {
			return nil
		}

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || i == attempts-1 {
			return err
		}

		// Ждем 2^i * 100 миллисекунд
		delay := time.Millisecond * 100 * (1 << i)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
