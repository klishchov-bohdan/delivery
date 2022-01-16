package main

import (
	store "github.com/klishchov-bohdan/delivery/pkg/store"
	"github.com/klishchov-bohdan/delivery/pkg/store/database/mysql"
	"github.com/klishchov-bohdan/delivery/pkg/store/database/transaction"
	"log"
)

func main() {
	db, err := mysql.Dial()
	if err != nil {
		log.Fatal(err)
	}
	storage := store.NewStore(mysql.NewUserMySQL(db))
	if err != nil {
		log.Fatal(err)
	}
	err = storage.ExecuteTransaction(db.GetDB(),
		transaction.NewPipelineStmt("insert into suppliers(id, name, description) values(?, ?, ?)", 1, "Supplier1", "desc1"),
		transaction.NewPipelineStmt("insert into suppliers(id, name, description) values(?, ?, ?)", 6, "Supplier6", "desc6"),
		transaction.NewPipelineStmt("INSERT INTO users(id, name, email, password_hash) VALUES(?, ?, ?, ?)", []byte("dscjhbkecn"), "Alex", "alex@gmail.com", "some hash"))
	if err != nil {
		log.Fatal(err)
	}
}
