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
	Type        string
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
		Type:        prod.Type,
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
	Type        string    `json:"type"`
}

func (mi *MenuItem) ToProduct(supplierID uuid.UUID) *Product {
	return &Product{
		ID:          mi.ID,
		SupplierID:  supplierID,
		Name:        mi.Name,
		Image:       mi.Image,
		Description: mi.Description,
		Type:        mi.Type,
		Price:       mi.Price,
		Weight:      mi.Weight,
	}
}
