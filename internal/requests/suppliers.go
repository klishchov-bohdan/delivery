package requests

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type SupplierWebRequest struct {
	Name        string                 `json:"name"`
	Image       string                 `json:"image"`
	Description string                 `json:"description"`
	WorkingTime models.WorkingSchedule `json:"working_time"`
	Menu        []*MenuItemRequest     `json:"menu"`
}

func (swr *SupplierWebRequest) GetSupplier() *models.Supplier {
	return &models.Supplier{
		ID:          uuid.New(),
		Name:        swr.Name,
		Image:       swr.Image,
		Description: swr.Description,
		WorkingTime: swr.WorkingTime,
	}
}

func (swr *SupplierWebRequest) GetProducts() *[]models.Product {
	var products []models.Product
	for _, menuItem := range swr.Menu {
		products = append(products, *menuItem.ToProduct(uuid.New()))
	}
	return &products
}

type MenuItemRequest struct {
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Weight      float64 `json:"weight"`
	Type        string  `json:"type"`
}

func (mig *MenuItemRequest) ToProduct(supplierID uuid.UUID) *models.Product {
	return &models.Product{
		ID:          uuid.New(),
		SupplierID:  supplierID,
		Name:        mig.Name,
		Image:       mig.Image,
		Description: mig.Description,
		Type:        mig.Type,
		Price:       mig.Price,
		Weight:      mig.Weight,
	}
}
