package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID
	SupplierID  uuid.UUID
	Name        string
	Image       string
	Description string
	Price       float64
	Weight      float64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (prod *Product) ToMenuItem() *MenuItem {
	return &MenuItem{
		ID:          prod.ID,
		Name:        prod.Name,
		Image:       prod.Image,
		Description: prod.Description,
		Price:       prod.Price,
		Weight:      prod.Weight,
	}
}

type MenuItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Weight      float64   `json:"weight"`
}
