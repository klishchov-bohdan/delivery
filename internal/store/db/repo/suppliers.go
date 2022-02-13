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
	rows, err := r.DB.Query("SELECT id, name, image, description, open_in, close_in, working_days, created_at, updated_at FROM suppliers")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var supplier models.Supplier
		err = rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Image,
			&supplier.Description,
			&supplier.WorkingTime.OpenIn,
			&supplier.WorkingTime.CloseIn,
			&supplier.WorkingTime.WorkingDays,
			&supplier.CreatedAt,
			&supplier.UpdatedAt)
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
		"SELECT id, name, image, description, open_in, close_in, working_days, created_at, updated_at FROM suppliers WHERE id = ?", uid).
		Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Image,
			&supplier.Description,
			&supplier.WorkingTime.OpenIn,
			&supplier.WorkingTime.CloseIn,
			&supplier.WorkingTime.WorkingDays,
			&supplier.CreatedAt,
			&supplier.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SuppliersRepo) CreateSupplier(supplier *models.Supplier) (uuid.UUID, error) {
	if supplier == nil {
		return uuid.Nil, errors.New("no supplier provided")
	}
	uid, err := supplier.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO suppliers(id, name, image, description, open_in, close_in, working_days) VALUES(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(
			uid,
			supplier.Name,
			supplier.Description,
			supplier.WorkingTime.OpenIn,
			supplier.WorkingTime.CloseIn,
			supplier.WorkingTime.WorkingDays)
		if err != nil {
			return uuid.Nil, err
		}
		return supplier.ID, nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO suppliers(id, name, image, description, open_in, close_in, working_days) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(
		uid,
		supplier.Name,
		supplier.Image,
		supplier.Description,
		supplier.WorkingTime.OpenIn,
		supplier.WorkingTime.CloseIn,
		supplier.WorkingTime.WorkingDays)
	if err != nil {
		return uuid.Nil, err
	}
	return supplier.ID, nil
}

func (r *SuppliersRepo) UpdateSupplier(supplier *models.Supplier) (uuid.UUID, error) {
	if supplier == nil {
		return uuid.Nil, errors.New("no supplier provided")
	}
	uid, err := supplier.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE suppliers SET name = ?, image = ?, description = ?, open_in = ?, close_in = ?, working_days = ? WHERE id = ?")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(
			supplier.Name,
			supplier.Image,
			supplier.Description,
			supplier.WorkingTime.OpenIn,
			supplier.WorkingTime.CloseIn,
			supplier.WorkingTime.WorkingDays,
			uid)
		if err != nil {
			return uuid.Nil, err
		}
		return supplier.ID, nil
	}
	stmt, err := r.DB.Prepare("UPDATE suppliers SET name = ?, image = ?, description = ?, open_in = ?, close_in = ?, working_days = ? WHERE id = ?")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(
		supplier.Name,
		supplier.Image,
		supplier.Description,
		supplier.WorkingTime.OpenIn,
		supplier.WorkingTime.CloseIn,
		supplier.WorkingTime.WorkingDays,
		uid)
	if err != nil {
		return uuid.Nil, err
	}
	return supplier.ID, nil
}

func (r *SuppliersRepo) DeleteSupplier(id uuid.UUID) (uuid.UUID, error) {
	uid, err := id.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM suppliers WHERE id = ?", uid)
		if err != nil {
			return uuid.Nil, err
		}
		return id, nil
	}
	_, err = r.DB.Exec("DELETE FROM suppliers WHERE id = ?", uid)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
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
