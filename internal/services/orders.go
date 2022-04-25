package services

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/responses"
	"github.com/klishchov-bohdan/delivery/internal/store"
)

type OrdersWebService struct {
	store *store.Store
}

func NewOrderWebService(store *store.Store) *OrdersWebService {
	return &OrdersWebService{
		store: store,
	}
}

func (service *OrdersWebService) GetOrderByID(id uuid.UUID) (*responses.OrderResponse, error) {
	order, err := service.store.Orders.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	address, err := service.store.Addresses.GetShippingAddressByID(order.ShippingAddressID)
	if err != nil {
		return nil, err
	}
	user, _ := service.store.Users.GetUserByID(order.UserID)
	return responses.NewOrderResponse(order, user, address), nil
}

func (service *OrdersWebService) GetOrdersByUserID(userID uuid.UUID) (*[]responses.OrderResponse, error) {
	orders, err := service.store.Orders.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}
	var ordersResponse []responses.OrderResponse
	for _, order := range *orders {
		address, err := service.store.Addresses.GetShippingAddressByID(order.ShippingAddressID)
		if err != nil {
			return nil, err
		}
		user, _ := service.store.Users.GetUserByID(order.UserID)
		ordersResponse = append(ordersResponse, *responses.NewOrderResponse(&order, user, address))
	}
	return &ordersResponse, nil
}

func (service *OrdersWebService) CreateOrder(order *models.Order, address *models.ShippingAddress) (uuid.UUID, error) {
	err := service.store.Orders.BeginTx()
	if err != nil {
		return uuid.Nil, err
	}
	service.store.Addresses.SetTx(service.store.Orders.GetTx())
	orderID, err := service.store.Orders.CreateOrder(order)
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	_, err = service.store.Addresses.CreateShippingAddress(address)
	if err != nil {
		_ = service.store.Addresses.RollbackTx()
		return uuid.Nil, err
	}
	err = service.store.Orders.CommitTx()
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	service.store.Orders.SetTx(nil)
	service.store.Addresses.SetTx(nil)
	return orderID, nil
}

func (service *OrdersWebService) UpdateOrder(order *models.Order, address *models.ShippingAddress) (uuid.UUID, error) {
	err := service.store.Orders.BeginTx()
	if err != nil {
		return uuid.Nil, err
	}
	service.store.Addresses.SetTx(service.store.Orders.GetTx())
	orderID, err := service.store.Orders.UpdateOrder(order)
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	_, err = service.store.Addresses.UpdateShippingAddress(address)
	if err != nil {
		_ = service.store.Addresses.RollbackTx()
		return uuid.Nil, err
	}
	err = service.store.Orders.CommitTx()
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	service.store.Orders.SetTx(nil)
	service.store.Addresses.SetTx(nil)
	return orderID, nil
}

func (service *OrdersWebService) DeleteOrder(id uuid.UUID) (uuid.UUID, error) {
	err := service.store.Orders.BeginTx()
	if err != nil {
		return uuid.Nil, err
	}
	service.store.Addresses.SetTx(service.store.Orders.GetTx())
	order, err := service.store.Orders.GetOrderByID(id)
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	orderID, err := service.store.Orders.DeleteOrder(order.ID)
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	_, err = service.store.Addresses.DeleteShippingAddress(order.ShippingAddressID)
	if err != nil {
		_ = service.store.Addresses.RollbackTx()
		return uuid.Nil, err
	}
	err = service.store.Orders.CommitTx()
	if err != nil {
		_ = service.store.Orders.RollbackTx()
		return uuid.Nil, err
	}
	service.store.Orders.SetTx(nil)
	service.store.Addresses.SetTx(nil)
	return orderID, nil
}
