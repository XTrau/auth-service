package domain

import (
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
