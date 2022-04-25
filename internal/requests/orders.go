package requests

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type OrderRequest struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	TotalPrice      uint64
	ClientPhone     string
	ShippingAddress *models.ShippingAddress
}

func (order *OrderRequest) GetOrder() *models.Order {
	return &models.Order{
		ID:                order.ID,
		UserID:            order.UserID,
		TotalPrice:        order.TotalPrice,
		ClientPhone:       order.ClientPhone,
		ShippingAddressID: order.ShippingAddress.ID,
	}
}
