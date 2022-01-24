package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	AccessHash  string
	RefreshHash string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type ResponseToken struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}
