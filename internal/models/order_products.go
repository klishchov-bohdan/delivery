package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderProducts struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	ProductID       uuid.UUID
	ProductQuantity uint64
	TotalPrice      uint64
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
