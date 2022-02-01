package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klishchov-bohdan/delivery/internal/routes"
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
	defer db.Close()
	store := store.NewStore(db)
	services, err := services.NewManager(store)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.GenerateRoutes(services, r)

	http.ListenAndServe(":8080", r)
}
