package models

import (
	"github.com/google/uuid"
	"time"
)

type ShippingAddress struct {
	ID        uuid.UUID
	ZIPCode   string
	Country   string
	County    string
	State     string
	City      string
	Street    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
