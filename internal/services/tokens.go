package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/store"
)

type TokensWebService struct {
	store *store.Store
}

func NewTokensWebService(store *store.Store) *TokensWebService {
	return &TokensWebService{
		store: store,
	}
}

func (service *TokensWebService) GetTokenByID(id uuid.UUID) (*models.Token, error) {
	token, err := service.store.Tokens.GetTokenByID(id)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("token with uuid %s not found", id.String())
	}
	return token, nil
}

func (service *TokensWebService) GetTokenByUserID(id uuid.UUID) (*models.Token, error) {
	token, err := service.store.Tokens.GetTokenByUserID(id)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("token with UserID %s not found", id.String())
	}
	return token, nil
}

func (service *TokensWebService) CreateToken(token *models.Token) (uuid.UUID, error) {
	id, err := service.store.Tokens.CreateToken(token)
	if err != nil {
		return uuid.Nil, err
	}
	created, _ := service.store.Tokens.GetTokenByID(id)
	if created == nil {
		return uuid.Nil, fmt.Errorf("can`t create token with uuid %s", id.String())
	}
	return id, nil
}

func (service *TokensWebService) UpdateToken(token *models.Token) (uuid.UUID, error) {
	id, err := service.store.Tokens.UpdateToken(token)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (service *TokensWebService) DeleteTokenByID(id uuid.UUID) (uuid.UUID, error) {
	id, err := service.store.Tokens.DeleteTokenByID(id)
	if err != nil {
		return uuid.Nil, err
	}
	deleted, _ := service.store.Tokens.GetTokenByID(id)
	if deleted != nil {
		return uuid.Nil, fmt.Errorf("can`t delete token with uuid %s", id.String())
	}
	return id, nil
}

func (service *TokensWebService) DeleteTokenByUserID(id uuid.UUID) (uuid.UUID, error) {
	id, err := service.store.Tokens.DeleteTokenByUserID(id)
	if err != nil {
		return uuid.Nil, err
	}
	deleted, _ := service.store.Tokens.GetTokenByUserID(id)
	if deleted != nil {
		return uuid.Nil, fmt.Errorf("can`t delete token with uuid %s", id.String())
	}
	return id, nil
}
