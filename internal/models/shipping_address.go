package models

import (
	"github.com/google/uuid"
	"time"
)

type ShippingAddress struct {
	ID        uuid.UUID
	ZIPCode   string
	Country   string
	Region    string
	Street    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
