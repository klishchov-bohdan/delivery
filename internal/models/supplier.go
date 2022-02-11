package models

import (
	"github.com/google/uuid"
	"time"
)

type Supplier struct {
	ID          uuid.UUID
	Name        string
	Image       string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type SupplierWeb struct {
	ID          uuid.UUID
	Name        string
	Image       string
	Description string
	Menu        []struct {
		ID          uuid.UUID
		Name        string
		Image       string
		Description string
		Price       float64
		Weight      float64
	}
}
