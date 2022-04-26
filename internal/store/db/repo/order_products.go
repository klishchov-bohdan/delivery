package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type OrderProductsRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewOrderProductsRepo(db *sql.DB) *OrderProductsRepo {
	return &OrderProductsRepo{DB: db}
}

func (r *OrderProductsRepo) GetProductsByOrderID(orderID uuid.UUID) (*[]models.OrderProducts, error) {
	var orderedProducts []models.OrderProducts
	uid, err := orderID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query("SELECT id, order_id, product_id, product_quantity, total_price, created_at, updated_at FROM order_products WHERE order_id = ?", uid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var orderProduct models.OrderProducts
		err = rows.Scan(&orderProduct.ID, &orderProduct.OrderID, &orderProduct.ProductID, &orderProduct.ProductQuantity, &orderProduct.TotalPrice, &orderProduct.CreatedAt, &orderProduct.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orderedProducts = append(orderedProducts, orderProduct)
	}
	return &orderedProducts, nil
}

func (r *OrderProductsRepo) CreateOrderProducts(orderProduct *models.OrderProducts) (uuid.UUID, error) {
	if orderProduct == nil {
		return uuid.Nil, errors.New("no ordered products provided")
	}
	uid, err := orderProduct.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	orderID, err := orderProduct.OrderID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	productID, err := orderProduct.ProductID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO order_products(id, order_id, product_id, product_quantity, total_price) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(uid, orderID, productID, orderProduct.ProductQuantity, orderProduct.TotalPrice)
		if err != nil {
			return uuid.Nil, err
		}
		return orderProduct.ID, nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO order_products(id, order_id, product_id, product_quantity, total_price) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(uid, orderID, productID, orderProduct.ProductQuantity, orderProduct.TotalPrice)
	if err != nil {
		return uuid.Nil, err
	}
	return orderProduct.ID, nil
}

func (r *OrderProductsRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *OrderProductsRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *OrderProductsRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *OrderProductsRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *OrderProductsRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}
