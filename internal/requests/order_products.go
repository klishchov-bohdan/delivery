package requests

import "github.com/google/uuid"

type OrderProductsRequest struct {
	ProductID       uuid.UUID
	ProductQuantity uint64
	TotalPrice      uint64
}
