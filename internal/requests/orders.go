package requests

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type OrderRequest struct {
	UserID          uuid.UUID
	Products        *[]OrderProductsRequest
	TotalPrice      uint64
	ClientPhone     string
	ShippingAddress *AddressRequest
}

func (order *OrderRequest) ToOrder() (*models.Order, *models.ShippingAddress, *[]models.OrderProducts) {
	newAddress := &models.ShippingAddress{
		ID:      uuid.New(),
		ZIPCode: order.ShippingAddress.ZIPCode,
		Country: order.ShippingAddress.Country,
		County:  order.ShippingAddress.County,
		State:   order.ShippingAddress.State,
		City:    order.ShippingAddress.City,
		Street:  order.ShippingAddress.Street,
	}
	newOrder := &models.Order{
		ID:                uuid.New(),
		UserID:            order.UserID,
		TotalPrice:        order.TotalPrice,
		ClientPhone:       order.ClientPhone,
		ShippingAddressID: newAddress.ID,
	}
	var newOrderedProducts []models.OrderProducts
	for _, orderProduct := range *order.Products {
		orderedProduct := &models.OrderProducts{
			ID:              uuid.New(),
			OrderID:         newOrder.ID,
			ProductID:       orderProduct.ProductID,
			ProductQuantity: orderProduct.ProductQuantity,
			TotalPrice:      orderProduct.TotalPrice,
		}
		newOrderedProducts = append(newOrderedProducts, *orderedProduct)
	}
	return newOrder, newAddress, &newOrderedProducts
}
