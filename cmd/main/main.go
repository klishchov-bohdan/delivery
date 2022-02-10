package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klishchov-bohdan/delivery/config"
	_ "github.com/klishchov-bohdan/delivery/docs"
	"github.com/klishchov-bohdan/delivery/internal/routes"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/store"
	"github.com/klishchov-bohdan/delivery/internal/store/db"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Go Restful API with Swagger
// @version 1.0
// @description Simple swagger implementation in Go HTTP

// @contact.email bogdan.bogdan2525@gmail.com

// @host localhost:8080
// @BasePath /
func main() {
	db, err := db.Dial()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	storage := store.NewStore(db)
	service, err := services.NewManager(storage)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.GenerateRoutes(service, cfg, r)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	http.ListenAndServe(":8080", r)
}
