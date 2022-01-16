package models

import "time"

type Product struct {
	ID          []byte
	Name        string
	Description string
	Price       float64
	Weight      float64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
