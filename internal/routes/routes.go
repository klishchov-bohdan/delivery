package routes

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/klishchov-bohdan/delivery/docs"
	"github.com/klishchov-bohdan/delivery/internal/controller"
	"github.com/klishchov-bohdan/delivery/internal/middleware"
	"github.com/klishchov-bohdan/delivery/internal/services"
)

func GenerateRoutes(services *services.Manager, r *chi.Mux) {
	ctr := controller.NewController(services)
	mw := middleware.NewMiddleware(services)
	r.Post("/login", ctr.Auth.Login)
	r.Post("/registration", ctr.Auth.Registration)
	r.Get("/refresh", ctr.Auth.Refresh)
	r.With(mw.AuthCheck).Get("/logout", ctr.Auth.Logout)

	r.With(mw.AuthCheck).Get("/profile", ctr.User.Profile)
}
