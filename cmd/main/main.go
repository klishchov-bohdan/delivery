package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/store/db"
	"github.com/klishchov-bohdan/delivery/internal/store/db/repo"
	"log"
)

func main() {
	db, err := db.Dial()
	if err != nil {
		log.Fatal(err)
	}
	ur := repo.NewUsersRepo(db)
	sr := repo.NewSuppliersRepo(db)
	pr := repo.NewProductMySQL(db)

	users, err := ur.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	suppliers, err := sr.GetAllSuppliers()
	if err != nil {
		log.Fatal(err)
	}

	products, err := pr.GetAllProducts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)
	fmt.Println(suppliers)
	fmt.Println(products)
	if err = ur.BeginTx(); err != nil {
		log.Fatal(err)
	}
	sr.SetTx(ur.GetTx())
	user, err := models.CreateUser("Luc", "luc@gmail.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	_, err = ur.CreateUser(user)
	if err != nil {
		log.Fatal(err)
	}
	err = sr.CreateSupplier(&models.Supplier{ID: uuid.New(), Name: "Name", Description: "Desc"})
	if err != nil {
		_ = ur.RollbackTx()
		log.Fatal(err)
	}
	err = ur.CommitTx()
	if err != nil {
		log.Fatal(err)
	}
	sr.SetTx(nil)
}
