package main

import (
	"github.com/klishchov-bohdan/delivery/internal/controller"
	"github.com/klishchov-bohdan/delivery/internal/middleware"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/store"
	"github.com/klishchov-bohdan/delivery/internal/store/db"
	"log"
	"net/http"
)

func main() {
	db, err := db.Dial()
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStore(db)
	services, err := services.NewManager(store)
	if err != nil {
		log.Fatal(err)
	}
	mw := middleware.NewMiddleware(services)
	tc := controller.NewTokenController(services)
	http.HandleFunc("/login", tc.Login)
	http.HandleFunc("/registration", tc.Registration)
	http.Handle("/logout", mw.AuthCheck(http.HandlerFunc(tc.Logout)))
	http.Handle("/profile", mw.AuthCheck(http.HandlerFunc(tc.Profile)))
	http.ListenAndServe(":8080", nil)
}
