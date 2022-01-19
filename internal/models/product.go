package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Weight      float64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
