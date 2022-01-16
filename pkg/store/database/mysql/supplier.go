package mysql

import (
	"errors"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type SupplierMySQL struct {
	db *MySQL
}

func NewSupplierMySQL(db *MySQL) *SupplierMySQL {
	return &SupplierMySQL{db: db}
}

func (supplierMySQL *SupplierMySQL) GetAllSuppliers() (*[]models.Supplier, error) {
	var suppliers []models.Supplier
	rows, err := supplierMySQL.db.Query("SELECT id, name, description, created_at, updated_at FROM suppliers")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var supplier models.Supplier
		err = rows.Scan(&supplier.ID, &supplier.Name, &supplier.Description, &supplier.CreatedAt, &supplier.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}
	return &suppliers, nil
}

func (supplierMySQL *SupplierMySQL) GetSupplierByID(id []byte) (*models.Supplier, error) {
	var supplier models.Supplier
	err := supplierMySQL.db.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM suppliers WHERE id = ?", id).
		Scan(&supplier.ID, &supplier.Name, &supplier.Description, &supplier.CreatedAt, &supplier.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (supplierMySQL *SupplierMySQL) CreateSupplier(supplier *models.Supplier) (err error) {
	if supplier == nil {
		return errors.New("no supplier provided")
	}
	stmt, err := supplierMySQL.db.Prepare("INSERT INTO Suppliers(id, name, description) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(supplier.ID, supplier.Name, supplier.Description)
	if err != nil {
		return err
	}
	return nil
}

func (supplierMySQL *SupplierMySQL) UpdateSupplier(supplier *models.Supplier) (err error) {
	if supplier == nil {
		return errors.New("no supplier provided")
	}
	err = supplierMySQL.db.QueryRow("SELECT * FROM suppliers WHERE id = ?", supplier.ID).Scan()
	if err != nil {
		return errors.New("supplier not found")
	}
	stmt, err := supplierMySQL.db.Prepare("UPDATE suppliers SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(supplier.Name, supplier.Description, supplier.ID)
	if err != nil {
		return err
	}
	return nil
}

func (supplierMySQL *SupplierMySQL) DeleteSupplier(id []byte) (err error) {
	err = supplierMySQL.db.QueryRow("SELECT * FROM suppliers WHERE id = ?", id).Scan()
	if err != nil {
		return errors.New("supplier not found")
	}
	_, err = supplierMySQL.db.Exec("DELETE FROM suppliers WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
