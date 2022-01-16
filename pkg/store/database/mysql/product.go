package mysql

import (
	"errors"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type ProductMySQL struct {
	db *MySQL
}

func NewProductMySQL(db *MySQL) *ProductMySQL {
	return &ProductMySQL{db: db}
}

func (productMySQL *ProductMySQL) GetAllProducts() (*[]models.Product, error) {
	var products []models.Product
	rows, err := productMySQL.db.Query("SELECT id, name, description, price, weight, created_at, updated_at FROM products")
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

func (productMySQL *ProductMySQL) GetProductByID(id []byte) (*models.Product, error) {
	var product models.Product
	err := productMySQL.db.QueryRow(
		"SELECT id, name, description, price, weight, created_at, updated_at FROM products WHERE id = ?", id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Weight, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (productMySQL *ProductMySQL) CreateProduct(product *models.Product) (err error) {
	if product == nil {
		return errors.New("no product provided")
	}
	stmt, err := productMySQL.db.Prepare("INSERT INTO products(id, name, description, price, weight) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ID, product.Name, product.Description, product.Price, product.Weight)
	if err != nil {
		return err
	}
	return nil
}

func (productMySQL *ProductMySQL) UpdateProduct(product *models.Product) (err error) {
	if product == nil {
		return errors.New("no product provided")
	}
	err = productMySQL.db.QueryRow("SELECT * FROM products WHERE id = ?", product.ID).Scan()
	if err != nil {
		return errors.New("product not found")
	}
	stmt, err := productMySQL.db.Prepare("UPDATE products SET name = ?, description = ?, price = ?, weight = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.Name, product.Description, product.Price, product.Weight, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (productMySQL *ProductMySQL) DeleteProduct(id []byte) (err error) {
	err = productMySQL.db.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan()
	if err != nil {
		return errors.New("product not found")
	}
	_, err = productMySQL.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
