package responses

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type OrderResponse struct {
	ID              uuid.UUID
	User            *models.User
	TotalPrice      uint64
	ClientPhone     string
	ShippingAddress *models.ShippingAddress
}

func NewOrderResponse(order *models.Order, user *models.User, address *models.ShippingAddress) *OrderResponse {
	return &OrderResponse{
		ID:              order.ID,
		User:            user,
		TotalPrice:      order.TotalPrice,
		ClientPhone:     order.ClientPhone,
		ShippingAddress: address,
	}
}
