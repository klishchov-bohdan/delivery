package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type ProductRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewProductMySQL(db *sql.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (r *ProductRepo) GetAllProducts() (*[]models.Product, error) {
	var products []models.Product
	rows, err := r.DB.Query("SELECT id, name, description, price, weight, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Weight, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return &products, nil
}

func (r *ProductRepo) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, name, description, price, weight, created_at, updated_at FROM products WHERE id = ?", uid).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Weight, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) CreateProduct(product *models.Product) (err error) {
	if product == nil {
		return errors.New("no product provided")
	}
	uid, err := product.ID.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO products(id, name, description, price, weight) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(uid, product.Name, product.Description, product.Price, product.Weight)
		if err != nil {
			return err
		}
		return nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO products(id, name, description, price, weight) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(uid, product.Name, product.Description, product.Price, product.Weight)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) UpdateProduct(product *models.Product) (err error) {
	if product == nil {
		return errors.New("no product provided")
	}
	err = r.DB.QueryRow("SELECT * FROM products WHERE id = ?", product.ID).Scan()
	if err != nil {
		return errors.New("product not found")
	}
	uid, err := product.ID.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE products SET name = ?, description = ?, price = ?, weight = ? WHERE id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(product.Name, product.Description, product.Price, product.Weight, uid)
		if err != nil {
			return err
		}
		return nil
	}
	stmt, err := r.DB.Prepare("UPDATE products SET name = ?, description = ?, price = ?, weight = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.Name, product.Description, product.Price, product.Weight, uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) DeleteProduct(id uuid.UUID) (err error) {
	err = r.DB.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan()
	if err != nil {
		return errors.New("product not found")
	}
	uid, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM products WHERE id = ?", uid)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = r.DB.Exec("DELETE FROM products WHERE id = ?", uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *ProductRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *ProductRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *ProductRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *ProductRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}
