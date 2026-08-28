package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/XTrau/auth-service/internal/domain"
	errs "github.com/XTrau/auth-service/internal/errors"
)

type PgUserRepository struct {
}

func NewPgUserRepository() *PgUserRepository {
	return &PgUserRepository{}
}

func (repo *PgUserRepository) Create(ctx context.Context, username string, passwordHash string) (*domain.User, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("Создание пользователя в базе данных", slog.String("username", username))
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"

	var id int
	err = tx.QueryRowContext(ctx, query, username, passwordHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           domain.UserID(id),
		Username:     username,
		PasswordHash: passwordHash,
	}, nil
}

func (repo *PgUserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("Получение пользователя с базы данных",
		slog.Int64("id", int64(id)),
	)
	query := "SELECT username, password_hash FROM users WHERE id=$1"

	var username, password string
	err = tx.QueryRowContext(ctx, query, id).Scan(&username, &password)

	if err == sql.ErrNoRows {
		return nil, errs.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           domain.UserID(id),
		Username:     username,
		PasswordHash: password,
	}, nil
}

func (repo *PgUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("Получение пользователя с базы данных",
		slog.String("username", username),
	)
	query := "SELECT id, password_hash FROM users WHERE username=$1"

	var id int64
	var password string
	err = tx.QueryRowContext(ctx, query, username).Scan(&id, &password)

	if err == sql.ErrNoRows {
		return nil, errs.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           domain.UserID(id),
		Username:     username,
		PasswordHash: password,
	}, nil
}
