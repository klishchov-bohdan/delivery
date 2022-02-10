package controller

import (
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/services"
)

type Controller struct {
	Auth *AuthController
	User *UserController
}

func NewController(services *services.Manager, cfg *config.Config) *Controller {
	return &Controller{
		Auth: NewAuthController(services),
		User: NewUserController(services, cfg),
	}
}
