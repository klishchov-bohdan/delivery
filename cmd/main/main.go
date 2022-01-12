package main

import (
	"fmt"
	"github.com/klishchov-bohdan/delivery/pkg/store/mysql"
	"log"
)

func main() {
	db, err := mysql.Dial()
	if err != nil {
		log.Fatal(err)
	}
	u := mysql.NewUserMySQL(db)
	users, err := u.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(users)
}
