package store

import (
	"database/sql"
	"github.com/klishchov-bohdan/delivery/internal/store/db/repo"
)

type Store struct {
	Users     *repo.UsersRepo
	Suppliers *repo.SuppliersRepo
	Products  *repo.ProductRepo
	Addresses *repo.ShippingAddressesRepo
	Tokens    *repo.TokensRepo
}

func NewStore(db *sql.DB) *Store {
	ur := repo.NewUsersRepo(db)
	sr := repo.NewSuppliersRepo(db)
	pr := repo.NewProductRepo(db)
	ar := repo.NewShippingAddressesRepo(db)
	tr := repo.NewTokensRepo(db)

	return &Store{
		Users:     ur,
		Suppliers: sr,
		Products:  pr,
		Addresses: ar,
		Tokens:    tr,
	}
}
