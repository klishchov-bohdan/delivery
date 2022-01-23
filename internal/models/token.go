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
