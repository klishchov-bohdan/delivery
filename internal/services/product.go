package services

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/store"
)

type ProductsWebService struct {
	store *store.Store
}

func NewProductsWebService(store *store.Store) *ProductsWebService {
	return &ProductsWebService{
		store: store,
	}
}

func (service *ProductsWebService) GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error) {
	product, err := service.store.Products.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return product.ToMenuItem(), nil
}

func (service *ProductsWebService) GetMenuBySupplierID(id uuid.UUID) (*[]models.MenuItem, error) {
	products, err := service.store.Products.GetProductsBySupplierID(id)
	if err != nil {
		return nil, err
	}
	var menu []models.MenuItem
	for _, product := range *products {
		menu = append(menu, *product.ToMenuItem())
	}
	return &menu, nil
}
