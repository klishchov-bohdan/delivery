package services

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/responses"
)

//go:generate mockgen --destination=mocks/users.go --package=mocks github.com/klishchov-bohdan/delivery/internal/services UserService
type UserService interface {
	GetAllUsers() (*[]models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) (uuid.UUID, error)
	UpdateUser(user *models.User) (uuid.UUID, error)
	DeleteUser(id uuid.UUID) (uuid.UUID, error)
	CreateUserWithTokens(user *models.User, token *models.Token) (userID uuid.UUID, err error)
}

//go:generate mockgen --destination=mocks/tokens.go --package=mocks github.com/klishchov-bohdan/delivery/internal/services TokenService
type TokenService interface {
	GetTokenByID(id uuid.UUID) (*models.Token, error)
	GetTokenByUserID(id uuid.UUID) (*models.Token, error)
	CreateToken(token *models.Token) (uuid.UUID, error)
	UpdateToken(token *models.Token) (uuid.UUID, error)
	DeleteTokenByID(id uuid.UUID) (uuid.UUID, error)
	DeleteTokenByUserID(id uuid.UUID) (uuid.UUID, error)
}

type SupplierService interface {
	GetSupplierByID(id uuid.UUID) (*models.SupplierWeb, error)
	GetAllSuppliers() (*[]models.SupplierWeb, error)
	CreateSupplier(supplierWeb *models.SupplierWeb) (uuid.UUID, error)
	UpdateSupplier(supplierWeb *models.SupplierWeb) (uuid.UUID, error)
	DeleteSupplier(id uuid.UUID) (uuid.UUID, error)
}

type ProductService interface {
	GetMenuBySupplierID(id uuid.UUID) (*[]models.MenuItem, error)
	GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error)
}

type OrderService interface {
	GetOrderByID(id uuid.UUID) (*responses.OrderResponse, error)
	GetOrdersByUserID(userID uuid.UUID) (*[]responses.OrderResponse, error)
	CreateOrder(order *models.Order, address *models.ShippingAddress, orderedProducts *[]models.OrderProducts) (uuid.UUID, error)
	UpdateOrder(order *models.Order, address *models.ShippingAddress) (uuid.UUID, error)
	DeleteOrder(id uuid.UUID) (uuid.UUID, error)
}
