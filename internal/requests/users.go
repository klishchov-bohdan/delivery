package requests

import "github.com/google/uuid"

type UserRegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
