package store

import (
	"database/sql"
	"github.com/klishchov-bohdan/delivery/pkg/store/database"
	"github.com/klishchov-bohdan/delivery/pkg/store/database/transaction"
)

type Store struct {
	Users     database.UsersDBInterface
	Suppliers database.SuppliersDBInterface
}

func NewStore(users database.UsersDBInterface) *Store {
	return &Store{
		Users: users,
	}
}

func (store *Store) ExecuteTransaction(db *sql.DB, stmts ...*transaction.PipelineStmt) error {
	err := transaction.WithTransaction(db, func(tx transaction.Transaction) error {
		_, err := transaction.RunPipeline(tx, stmts...)
		return err
	})
	return err
}
