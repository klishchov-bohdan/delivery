package services

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type UserService interface {
	GetAllUsers() (*[]models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) (uuid.UUID, error)
	UpdateUser(user *models.User) (uuid.UUID, error)
	DeleteUser(id uuid.UUID) (uuid.UUID, error)
}

type TokenService interface {
	GetTokenByID(id uuid.UUID) (*models.Token, error)
	GetTokenByUserID(id uuid.UUID) (*models.Token, error)
	CreateToken(token *models.Token) (uuid.UUID, error)
	UpdateToken(token *models.Token) (uuid.UUID, error)
	DeleteTokenByID(id uuid.UUID) (uuid.UUID, error)
	DeleteTokenByUserID(id uuid.UUID) (uuid.UUID, error)
}
