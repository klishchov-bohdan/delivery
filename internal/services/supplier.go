package services

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/store"
)

type SuppliersWebService struct {
	store *store.Store
}

func NewSuppliersWebService(store *store.Store) *SuppliersWebService {
	return &SuppliersWebService{
		store: store,
	}
}

func (service *SuppliersWebService) GetAllSuppliers() (*[]models.SupplierWeb, error) {
	suppliers, err := service.store.Suppliers.GetAllSuppliers()
	if err != nil {
		return nil, err
	}
	var suppliersWeb []models.SupplierWeb
	for _, supplier := range *suppliers {
		products, err := service.store.Products.GetProductsBySupplierID(supplier.ID)
		if err != nil {
			return nil, err
		}
		suppliersWeb = append(suppliersWeb, *models.NewSupplierWeb(&supplier, products))
	}
	return &suppliersWeb, nil
}

func (service *SuppliersWebService) GetSupplierByID(id uuid.UUID) (*models.SupplierWeb, error) {
	supplier, err := service.store.Suppliers.GetSupplierByID(id)
	if err != nil {
		return nil, err
	}
	products, err := service.store.Products.GetProductsBySupplierID(supplier.ID)
	if err != nil {
		return nil, err
	}
	supplierWeb := models.NewSupplierWeb(supplier, products)
	return supplierWeb, nil
}

func (service *SuppliersWebService) CreateSupplier(supplierWeb *models.SupplierWeb) (uuid.UUID, error) {
	err := service.store.Suppliers.BeginTx()
	if err != nil {
		return uuid.Nil, err
	}
	service.store.Products.SetTx(service.store.Suppliers.GetTx())
	supplierID, err := service.store.Suppliers.CreateSupplier(&models.Supplier{
		ID:          supplierWeb.ID,
		Name:        supplierWeb.Name,
		Image:       supplierWeb.Image,
		Description: supplierWeb.Description,
		WorkingTime: supplierWeb.WorkingTime,
	})
	if err != nil {
		_ = service.store.Suppliers.RollbackTx()
		return uuid.Nil, err
	}
	for _, menuItem := range supplierWeb.Menu {
		_, err = service.store.Products.CreateProduct(&models.Product{
			ID:          menuItem.ID,
			SupplierID:  supplierID,
			Name:        menuItem.Name,
			Image:       menuItem.Image,
			Description: menuItem.Description,
			Price:       menuItem.Price,
			Weight:      menuItem.Weight,
		})
		if err != nil {
			_ = service.store.Products.RollbackTx()
			return uuid.Nil, err
		}
	}
	err = service.store.Suppliers.CommitTx()
	if err != nil {
		_ = service.store.Suppliers.RollbackTx()
		return uuid.Nil, err
	}
	service.store.Suppliers.SetTx(nil)
	service.store.Products.SetTx(nil)
	return supplierID, nil
}

func (service *SuppliersWebService) UpdateSupplier(supplierWeb *models.SupplierWeb) (uuid.UUID, error) {
	err := service.store.Suppliers.BeginTx()
	if err != nil {
		return uuid.Nil, err
	}
	service.store.Products.SetTx(service.store.Suppliers.GetTx())
	supplierID, err := service.store.Suppliers.UpdateSupplier(&models.Supplier{
		ID:          supplierWeb.ID,
		Name:        supplierWeb.Name,
		Image:       supplierWeb.Image,
		Description: supplierWeb.Description,
		WorkingTime: supplierWeb.WorkingTime,
	})
	if err != nil {
		_ = service.store.Suppliers.RollbackTx()
		return uuid.Nil, err
	}
	for _, menuItem := range supplierWeb.Menu {
		_, err = service.store.Products.UpdateProduct(&models.Product{
			ID:          menuItem.ID,
			SupplierID:  supplierID,
			Name:        menuItem.Name,
			Image:       menuItem.Image,
			Description: menuItem.Description,
			Price:       menuItem.Price,
			Weight:      menuItem.Weight,
		})
		if err != nil {
			_ = service.store.Products.RollbackTx()
			return uuid.Nil, err
		}
	}
	err = service.store.Suppliers.CommitTx()
	if err != nil {
		_ = service.store.Suppliers.RollbackTx()
		return uuid.Nil, err
	}
	service.store.Suppliers.SetTx(nil)
	service.store.Products.SetTx(nil)
	return supplierID, nil
}

func (service *SuppliersWebService) DeleteSupplier(id uuid.UUID) (uuid.UUID, error) {
	err := service.store.Suppliers.BeginTx()
	if err != nil {
		return uuid.Nil, err
	}
	service.store.Products.SetTx(service.store.Suppliers.GetTx())
	supplierID, err := service.store.Suppliers.DeleteSupplier(id)
	if err != nil {
		_ = service.store.Suppliers.RollbackTx()
		return uuid.Nil, err
	}
	_, err = service.store.Products.DeleteProductsBySupplierID(supplierID)
	if err != nil {
		_ = service.store.Products.RollbackTx()
		return uuid.Nil, err
	}
	err = service.store.Suppliers.CommitTx()
	if err != nil {
		_ = service.store.Suppliers.RollbackTx()
		return uuid.Nil, err
	}
	service.store.Suppliers.SetTx(nil)
	service.store.Products.SetTx(nil)
	return supplierID, nil
}
