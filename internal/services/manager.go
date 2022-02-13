package services

import (
	"errors"
	"github.com/klishchov-bohdan/delivery/internal/store"
)

type Manager struct {
	User     UserService
	Token    TokenService
	Supplier SupplierService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("no store provided")
	}
	return &Manager{
		User:     NewUsersWebService(store),
		Token:    NewTokensWebService(store),
		Supplier: NewSuppliersWebService(store),
	}, nil
}
