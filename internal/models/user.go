package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           []byte
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func CreateUser(name, email, password string) (*User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	id, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(pwdHash),
	}, nil
}
