package controller

import (
	"encoding/json"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"net/http"
	"os"
)

type UserController struct {
	services *services.Manager
}

func NewUserController(services *services.Manager) *UserController {
	return &UserController{
		services: services,
	}
}

func (ctr *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		accessString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		accessSecret := os.Getenv("AccessSecret")
		claims, err := token.GetClaims(accessString, accessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		user, err := ctr.services.User.GetUserByID(claims.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}
