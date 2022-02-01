package controller

import "github.com/klishchov-bohdan/delivery/internal/services"

type Controller struct {
	Auth *AuthController
	User *UserController
}

func NewController(services *services.Manager) *Controller {
	return &Controller{
		Auth: NewAuthController(services),
		User: NewUserController(services),
	}
}
