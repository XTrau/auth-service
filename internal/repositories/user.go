package repositories

import (
	"database/sql"
	"log/slog"

	"github.com/XTrau/auth-service/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (repo *PostgresUserRepository) Create(username string, passwordHash string) (*domain.User, error) {
	slog.Info("Создание пользователя в базе данных", slog.String("username", username))
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"

	var id int
	err := repo.db.QueryRow(query, username, passwordHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           domain.UserID(id),
		Username:     username,
		PasswordHash: passwordHash,
	}, nil
}

func (repo *PostgresUserRepository) GetByID(id domain.UserID) (*domain.User, error) {
	slog.Info("Получение пользователя с базы данных",
		slog.Int64("id", int64(id)),
	)
	query := "SELECT username, password_hash FROM users WHERE id=$1"

	var username, password string
	err := repo.db.QueryRow(query, id).Scan(&username, &password)

	if err == sql.ErrNoRows {
		return nil, nil
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

func (repo *PostgresUserRepository) GetByUsername(username string) (*domain.User, error) {
	slog.Info("Получение пользователя с базы данных",
		slog.String("username", username),
	)
	query := "SELECT id, password_hash FROM users WHERE username=$1"

	var id int64
	var password string
	err := repo.db.QueryRow(query, username).Scan(&id, &password)

	if err == sql.ErrNoRows {
		return nil, nil
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
