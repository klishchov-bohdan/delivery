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

// GetProfile godoc
// @Summary Profile
// @Description Profile
// @Tags user
// @ID user-profile
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization"
// @Success 200 {object} models.User
// @Security ApiKeyAuth
// @Router /profile [get]
func (ctr *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
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
