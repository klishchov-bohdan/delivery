package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type ShippingAddressesRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewShippingAddressesRepo(db *sql.DB) *ShippingAddressesRepo {
	return &ShippingAddressesRepo{DB: db}
}

func (r *ShippingAddressesRepo) GetAllShippingAddresses() (*[]models.ShippingAddress, error) {
	var addresses []models.ShippingAddress
	rows, err := r.DB.Query("SELECT id, zip_code, country, state, city, street, county, created_at, updated_at FROM shipping_addresses")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var address models.ShippingAddress
		err = rows.Scan(&address.ID, &address.ZIPCode, &address.Country, &address.State, &address.City, &address.Street, &address.County, &address.CreatedAt, &address.UpdatedAt)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return &addresses, nil
}

func (r *ShippingAddressesRepo) GetShippingAddressByID(id uuid.UUID) (*models.ShippingAddress, error) {
	var address models.ShippingAddress
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, zip_code, country, state, city, street, county, created_at, updated_at FROM shipping_addresses WHERE id = ?", uid).
		Scan(&address.ID, &address.ZIPCode, &address.Country, &address.State, &address.City, &address.Street, &address.County, &address.CreatedAt, &address.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *ShippingAddressesRepo) CreateShippingAddress(address *models.ShippingAddress) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, errors.New("no address provided")
	}
	uid, err := address.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO shipping_addresses(id, zip_code, country, state, city, street, county) VALUES(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(uid, address.ZIPCode, address.Country, address.State, address.City, address.Street, address.County)
		if err != nil {
			return uuid.Nil, err
		}
		return address.ID, nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO shipping_addresses(id, zip_code, country, state, city, street, county) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(uid, address.ZIPCode, address.Country, address.State, address.City, address.Street, address.County)
	if err != nil {
		return uuid.Nil, err
	}
	return address.ID, nil
}

func (r *ShippingAddressesRepo) UpdateShippingAddress(address *models.ShippingAddress) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, errors.New("no address provided")
	}
	uid, err := address.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE shipping_addresses SET zip_code = ?, country = ?, state = ?, city = ?, street = ?, county = ? WHERE id = ?")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(address.ZIPCode, address.Country, address.State, address.City, address.Street, address.County, uid)
		if err != nil {
			return uuid.Nil, err
		}
		return address.ID, nil
	}
	stmt, err := r.DB.Prepare("UPDATE shipping_addresses SET zip_code = ?, country = ?, state = ?, city = ?, street = ?, county = ? WHERE id = ?")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(address.ZIPCode, address.Country, address.State, address.City, address.Street, address.County, uid)
	if err != nil {
		return uuid.Nil, err
	}
	return address.ID, nil
}

func (r *ShippingAddressesRepo) DeleteShippingAddress(id uuid.UUID) (uuid.UUID, error) {
	uid, err := id.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM shipping_addresses WHERE id = ?", uid)
		if err != nil {
			return uuid.Nil, err
		}
		return id, nil
	}
	_, err = r.DB.Exec("DELETE FROM shipping_addresses WHERE id = ?", uid)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *ShippingAddressesRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *ShippingAddressesRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *ShippingAddressesRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *ShippingAddressesRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *ShippingAddressesRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}
