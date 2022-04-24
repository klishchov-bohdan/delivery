package repo

import "database/sql"

type OrderRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (r *OrderRepo) GetOrdersByUserID() *models
