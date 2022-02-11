package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/klishchov-bohdan/delivery/config"
	_ "github.com/klishchov-bohdan/delivery/docs"
	"github.com/klishchov-bohdan/delivery/internal/controller"
	"github.com/klishchov-bohdan/delivery/internal/middleware"
	"github.com/klishchov-bohdan/delivery/internal/services"
)

func GenerateRoutes(services *services.Manager, cfg *config.Config, r *chi.Mux) {
	ctr := controller.NewController(services, cfg)
	mw := middleware.NewMiddleware(services)
	r.Post("/login", ctr.Auth.Login)
	r.Post("/registration", ctr.Auth.Registration)
	r.Post("/refresh", ctr.Auth.Refresh)
	r.With(mw.AuthCheck).Post("/logout", ctr.Auth.Logout)

	r.With(mw.AuthCheck).Get("/profile", ctr.User.GetProfile)
	r.With(mw.AuthCheck).Put("/profile", ctr.User.UpdateProfile)
	r.Get("/users/{userId}", ctr.User.GetUser)
	r.Get("/users", ctr.User.GetAllUsers)
	r.Post("/users", ctr.User.CreateUser)
	r.Put("/users", ctr.User.UpdateUser)
	r.Delete("/users/{userId}", ctr.User.DeleteUser)
}
