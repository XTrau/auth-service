package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type PgUnitOfWork struct {
	db *sql.DB
}

func NewPgUnitOfWork(db *sql.DB) *PgUnitOfWork {
	return &PgUnitOfWork{db}
}

func (uow *PgUnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(WithTx(ctx, tx)); err != nil {
		return err
	}

	return tx.Commit()
}

func (uow *PgUnitOfWork) ExecuteWithRetry(ctx context.Context, attempts int, fn func(ctx context.Context) error) error {
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
