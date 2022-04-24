package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	TotalPrice        uint64
	ClientPhone       string
	ShippingAddressID uuid.UUID
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}
