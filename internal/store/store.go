package store

import (
	"database/sql"
	"github.com/klishchov-bohdan/delivery/internal/store/db"
	transaction2 "github.com/klishchov-bohdan/delivery/internal/store/db/transaction"
)

type Store struct {
	Users     db.UsersDBInterface
	Suppliers db.SuppliersDBInterface
}

func NewStore(users db.UsersDBInterface) *Store {
	return &Store{
		Users: users,
	}
}

func (store *Store) ExecuteTransaction(db *sql.DB, stmts ...*transaction2.PipelineStmt) error {
	err := transaction2.WithTransaction(db, func(tx transaction2.Transaction) error {
		_, err := transaction2.RunPipeline(tx, stmts...)
		return err
	})
	return err
}
