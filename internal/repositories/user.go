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

func (repo *PgUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	slog.Info("Создание пользователя в базе данных", slog.String("username", user.Username))
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"

	var id int64
	err = tx.QueryRowContext(ctx, query, user.Username, user.PasswordHash).Scan(&id)
	if err != nil {
		return domain.User{}, err
	}

	user.ID = id
	return user, nil
}

func (repo *PgUserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	slog.Info("Получение пользователя с базы данных", slog.Int64("id", int64(id)))
	query := "SELECT username, password_hash FROM users WHERE id=$1"

	user := domain.User{ID: id}
	err = tx.QueryRowContext(ctx, query, id).Scan(&user.Username, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return domain.User{}, errs.ErrUserNotFound
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repo *PgUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	slog.Info("Получение пользователя с базы данных", slog.String("username", username))
	query := "SELECT id, password_hash FROM users WHERE username=$1"

	var user domain.User
	err = tx.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return domain.User{}, errs.ErrUserNotFound
	}

	if err != nil {
		return domain.User{}, err
	}

	user.Username = username
	return user, nil
}
