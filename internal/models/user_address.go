package models

import (
	"github.com/google/uuid"
	"time"
)

type UserAddress struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ZIPCode   string
	Country   string
	Region    string
	Street    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
