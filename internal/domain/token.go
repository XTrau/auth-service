package domain

import (
	"time"
)

type TokenPayload struct {
	Subject   int64  `json:"id"`
	Username  string `json:"username"`
	ExpiresAt time.Time
}

type TokenPair struct {
	Access  string
	Refresh string
}
