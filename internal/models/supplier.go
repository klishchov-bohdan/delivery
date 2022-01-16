package models

import "time"

type Supplier struct {
	ID          []byte
	Name        string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
