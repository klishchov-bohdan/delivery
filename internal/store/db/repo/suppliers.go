package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type SuppliersRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewSuppliersRepo(db *sql.DB) *SuppliersRepo {
	return &SuppliersRepo{DB: db}
}

func (r *SuppliersRepo) GetAllSuppliers() (*[]models.Supplier, error) {
	var suppliers []models.Supplier
	rows, err := r.DB.Query("SELECT id, name, description, created_at, updated_at FROM suppliers")
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

func (r *SuppliersRepo) GetSupplierByID(id uuid.UUID) (*models.Supplier, error) {
	var supplier models.Supplier
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM suppliers WHERE id = ?", uid).
		Scan(&supplier.ID, &supplier.Name, &supplier.Description, &supplier.CreatedAt, &supplier.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SuppliersRepo) CreateSupplier(supplier *models.Supplier) (err error) {
	if supplier == nil {
		return errors.New("no supplier provided")
	}
	uid, err := supplier.ID.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO suppliers(id, name, description) VALUES(?, ?, ?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(uid, supplier.Name, supplier.Description)
		if err != nil {
			return err
		}
		return nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO suppliers(id, name, description) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(uid, supplier.Name, supplier.Description)
	if err != nil {
		return err
	}
	return nil
}

func (r *SuppliersRepo) UpdateSupplier(supplier *models.Supplier) (err error) {
	if supplier == nil {
		return errors.New("no supplier provided")
	}
	err = r.DB.QueryRow("SELECT * FROM suppliers WHERE id = ?", supplier.ID).Scan()
	if err != nil {
		return errors.New("supplier not found")
	}
	uid, err := supplier.ID.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE suppliers SET name = ?, description = ? WHERE id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(supplier.Name, supplier.Description, uid)
		if err != nil {
			return err
		}
		return nil
	}
	stmt, err := r.DB.Prepare("UPDATE suppliers SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(supplier.Name, supplier.Description, uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *SuppliersRepo) DeleteSupplier(id uuid.UUID) (err error) {
	err = r.DB.QueryRow("SELECT * FROM suppliers WHERE id = ?", id).Scan()
	if err != nil {
		return errors.New("supplier not found")
	}
	uid, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM suppliers WHERE id = ?", uid)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = r.DB.Exec("DELETE FROM suppliers WHERE id = ?", uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *SuppliersRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *SuppliersRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *SuppliersRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *SuppliersRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *SuppliersRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}