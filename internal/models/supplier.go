package models

import (
	"github.com/google/uuid"
	"time"
)

type Supplier struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
