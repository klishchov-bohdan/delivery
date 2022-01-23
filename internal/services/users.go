package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/store"
)

type UsersWebService struct {
	store *store.Store
}

func NewUsersWebService(store *store.Store) *UsersWebService {
	return &UsersWebService{
		store: store,
	}
}

func (service *UsersWebService) GetAllUsers() (*[]models.User, error) {
	users, err := service.store.Users.GetAllUsers()
	if err != nil {
		return nil, err
	}
	if users == nil {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (service *UsersWebService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := service.store.Users.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with uuid %s not found", id.String())
	}
	return user, nil
}

func (service *UsersWebService) GetUserByEmail(email string) (*models.User, error) {
	user, err := service.store.Users.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	return user, nil
}

func (service *UsersWebService) CreateUser(user *models.User) (uuid.UUID, error) {
	id, err := service.store.Users.CreateUser(user)
	if err != nil {
		return uuid.Nil, err
	}
	created, _ := service.store.Users.GetUserByID(id)
	if created == nil {
		return uuid.Nil, fmt.Errorf("can`t create user with uuid %s", id.String())
	}
	return id, nil
}

func (service *UsersWebService) UpdateUser(user *models.User) (uuid.UUID, error) {
	id, err := service.store.Users.UpdateUser(user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (service *UsersWebService) DeleteUser(id uuid.UUID) (uuid.UUID, error) {
	id, err := service.store.Users.DeleteUser(id)
	if err != nil {
		return uuid.Nil, err
	}
	deleted, _ := service.store.Users.GetUserByID(id)
	if deleted != nil {
		return uuid.Nil, fmt.Errorf("can`t delete user with uuid %s", id.String())
	}
	return id, nil
}
