package database

import (
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type UsersDBInterface interface {
	GetAllUsers() (*[]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id []byte) (*models.User, error)
	CreateUser(user *models.User) (err error)
	UpdateUser(user *models.User) (err error)
	DeleteUser(id []byte) (err error)
}

type SuppliersDBInterface interface {
	GetAllSuppliers() (*[]models.Supplier, error)
	GetSupplierByID(id []byte) (*models.Supplier, error)
	CreateSupplier(user *models.Supplier) (err error)
	UpdateSupplier(user *models.Supplier) (err error)
	DeleteSupplier(id []byte) (err error)
}

type ProductsDBInterface interface {
	GetAllProducts() (*[]models.Product, error)
	GetProductByID(id []byte) (*models.Product, error)
	CreateProduct(user *models.Product) (err error)
	UpdateProduct(user *models.Product) (err error)
	DeleteProduct(id []byte) (err error)
}
