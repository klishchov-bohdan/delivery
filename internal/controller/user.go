package controller

import (
	"encoding/json"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"net/http"
)

type UserController struct {
	services *services.Manager
	cfg      *config.Config
}

func NewUserController(services *services.Manager, cfg *config.Config) *UserController {
	return &UserController{
		services: services,
		cfg:      cfg,
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
		claims, err := token.GetClaims(accessString, ctr.cfg.AccessSecret)
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
