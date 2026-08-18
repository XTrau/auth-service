package domain

import (
	"context"
	"strconv"
)

type UserID int64

func (u UserID) String() string {
	return strconv.FormatInt(int64(u), 10)
}

type UserCreate struct {
	Username string
	Password string
}

type User struct {
	ID           UserID
	Username     string
	PasswordHash string
}

type UserRepository interface {
	Create(ctx context.Context, username string, passwordHash string) (*User, error)
	GetByID(ctx context.Context, id UserID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
