package controller

import (
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/services"
)

type Controller struct {
	Auth     *AuthController
	User     *UserController
	Supplier *SupplierController
	Order    *OrderController
}

func NewController(services *services.Manager, cfg *config.Config) *Controller {
	return &Controller{
		Auth:     NewAuthController(services),
		User:     NewUserController(services, cfg),
		Supplier: NewSupplierController(services, cfg),
		Order:    NewOrderController(services, cfg),
	}
}
