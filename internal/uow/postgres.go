package uow

import (
	"context"
	"database/sql"

	"github.com/XTrau/auth-service/internal/domain"
	"github.com/XTrau/auth-service/internal/repositories"
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
		return nil
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
