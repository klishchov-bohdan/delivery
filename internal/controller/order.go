package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/requests"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"io/ioutil"
	"net/http"
)

type OrderController struct {
	services *services.Manager
	cfg      *config.Config
}

func NewOrderController(services *services.Manager, cfg *config.Config) *OrderController {
	return &OrderController{
		services: services,
		cfg:      cfg,
	}
}

// GetOrder godoc
// @Summary Get Order
// @Description Get Order
// @Tags order
// @ID get-order
// @Accept  json
// @Produce  json
// @Param orderId path string true "Order ID"
// @Success 200 {object} responses.OrderResponse
// @Router /orders/{orderId} [get]
func (ctr *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := chi.URLParam(r, "orderId")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order, err := ctr.services.Order.GetOrderByID(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// GetOrdersByUserID godoc
// @Summary Get Orders By User ID
// @Description Get Orders By User ID
// @Tags orders
// @ID get-orders-by-user-id
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {array} responses.OrderResponse
// @Router /orders/byUser/{userId} [get]
func (ctr *OrderController) GetOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := chi.URLParam(r, "userId")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		orders, err := ctr.services.Order.GetOrdersByUserID(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orders)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// CreateOrder godoc
// @Summary Create Order
// @Description Create Order
// @Tags orders
// @ID create-order
// @Accept  json
// @Produce  json
// @Param user body requests.OrderRequest true "Create Order Input"
// @Success 200 {object} responses.OrderResponse
// @Router /orders [post]
func (ctr *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var orderReq requests.OrderRequest
		err := json.Unmarshal(body, &orderReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdOrderID, err := ctr.services.Order.CreateOrder(
			orderReq.GetOrder(),
			orderReq.ShippingAddress)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		createdOrder, _ := ctr.services.Order.GetOrderByID(createdOrderID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdOrder)
	default:
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
	}
}
