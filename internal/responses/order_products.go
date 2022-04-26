package responses

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type OrderProductsResponse struct {
	ID              uuid.UUID
	Product         *models.Product
	ProductQuantity uint64
	TotalPrice      uint64
}
