package main

import (
	"fmt"
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

	users, err := ur.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)
	user, err := ur.GetUserByEmail("vlad@gmail.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	user.Name = "vlad123"
	user.Email = "vlad321@gmail.com"
	err = ur.UpdateUser(user)
	if err != nil {
		log.Fatal(err)
	}
}
