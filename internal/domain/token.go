package domain

import (
	"time"
)

type TokenPayload struct {
	Subject   UserID `json:"id"`
	Username  string `json:"username"`
	ExpiresAt time.Time
}

type TokenPair struct {
	Access  string
	Refresh string
}
