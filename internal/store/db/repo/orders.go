package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type OrderRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (r *OrderRepo) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, user_id, total_price, client_phone, shipping_address_id, created_at, updated_at FROM orders WHERE id = ?", uid).
		Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.ClientPhone, &order.ShippingAddressID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) GetOrdersByUserID(userID uuid.UUID) (*[]models.Order, error) {
	var orders []models.Order
	uid, err := userID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query("SELECT id, user_id, total_price, client_phone, shipping_address_id, created_at, updated_at FROM orders WHERE user_id = ?", uid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.ClientPhone, &order.ShippingAddressID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return &orders, nil
}

func (r *OrderRepo) CreateOrder(order *models.Order) (uuid.UUID, error) {
	if order == nil {
		return uuid.Nil, errors.New("no order provided")
	}
	uid, err := order.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := order.UserID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	shippingAddressID, err := order.ShippingAddressID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO orders(id, user_id, total_price, client_phone, shipping_address_id) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(uid, userID, order.TotalPrice, order.ClientPhone, shippingAddressID)
		if err != nil {
			return uuid.Nil, err
		}
		return order.ID, nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO orders(id, user_id, total_price, client_phone, shipping_address_id) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(uid, userID, order.TotalPrice, order.ClientPhone, shippingAddressID)
	if err != nil {
		return uuid.Nil, err
	}
	return order.ID, nil
}

func (r *OrderRepo) UpdateOrder(order *models.Order) (uuid.UUID, error) {
	if order == nil {
		return uuid.Nil, errors.New("no order provided")
	}
	uid, err := order.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := order.UserID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	shippingAddressID, err := order.ShippingAddressID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE orders SET user_id = ?, total_price = ?, client_phone = ?, shipping_address_id = ? WHERE id = ?")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(userID, order.TotalPrice, order.ClientPhone, shippingAddressID, uid)
		if err != nil {
			return uuid.Nil, err
		}
		return order.ID, nil
	}
	stmt, err := r.DB.Prepare("UPDATE orders SET user_id = ?, total_price = ?, client_phone = ?, shipping_address_id = ? WHERE id = ?")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(userID, order.TotalPrice, order.ClientPhone, shippingAddressID, uid)
	if err != nil {
		return uuid.Nil, err
	}
	return order.ID, nil
}

func (r *OrderRepo) DeleteOrder(id uuid.UUID) (uuid.UUID, error) {
	uid, err := id.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM orders WHERE id = ?", uid)
		if err != nil {
			return uuid.Nil, err
		}
		return id, nil
	}
	_, err = r.DB.Exec("DELETE FROM orders WHERE id = ?", uid)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *OrderRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *OrderRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *OrderRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *OrderRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *OrderRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}
