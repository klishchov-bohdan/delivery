package models

type UserDB struct {
	ID           uint64
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}
