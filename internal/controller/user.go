package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"io/ioutil"
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
// @Tags users
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

// UpdateProfile godoc
// @Summary Update Profile
// @Description Update Profile
// @Tags users
// @ID user-update-profile
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization"
// @Param user body models.UserToUpdate true "Update Profile Input"
// @Success 200 {object} models.User
// @Security ApiKeyAuth
// @Router /profile [put]
func (ctr *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		accessString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		claims, err := token.GetClaims(accessString, ctr.cfg.AccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		_, err = ctr.services.User.GetUserByID(claims.ID)
		if err != nil {
			http.Error(w, "User does not exists", http.StatusUnauthorized)
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var userReq models.UserToUpdate
		err = json.Unmarshal(body, &userReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if userReq.Name == "" || userReq.Email == "" || userReq.Password == "" {
			http.Error(w, "Cant be empty values", http.StatusBadRequest)
			return
		}
		userToUpdate, err := models.CreateUser(userReq.Name, userReq.Email, userReq.Password)
		userToUpdate.ID = userReq.ID
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updatedUserID, err := ctr.services.User.UpdateUser(userToUpdate)
		if err != nil || updatedUserID != userToUpdate.ID {
			http.Error(w, "Cant update profile", http.StatusInternalServerError)
			return
		}
		updatedUser, _ := ctr.services.User.GetUserByID(updatedUserID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)

	default:
		http.Error(w, "Only PUT method", http.StatusMethodNotAllowed)
	}
}

// GetUser godoc
// @Summary Get User
// @Description Get User
// @Tags users
// @ID get-user
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{userId} [get]
func (ctr *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := chi.URLParam(r, "userId")
		user, err := ctr.services.User.GetUserByID(uuid.MustParse(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// GetAllUsers godoc
// @Summary Get Users
// @Description Get Users
// @Tags users
// @ID get-users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func (ctr *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, err := ctr.services.User.GetAllUsers()
		if err != nil {
			http.Error(w, "Cant get all users", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// CreateUser godoc
// @Summary Create User
// @Description Create User
// @Tags users
// @ID create-user
// @Accept  json
// @Produce  json
// @Param user body models.User true "Create User Input"
// @Success 200 {object} models.User
// @Router /users [post]
func (ctr *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var userReq models.UserRegistrationRequest
		err := json.Unmarshal(body, &userReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, _ := models.CreateUser(userReq.Name, userReq.Email, userReq.Password)
		createdUserID, err := ctr.services.User.CreateUser(user)
		if err != nil || createdUserID != user.ID {
			http.Error(w, "Cant create user", http.StatusInternalServerError)
			return
		}
		createdUser, _ := ctr.services.User.GetUserByID(createdUserID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)
	default:
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
	}
}

// UpdateUser godoc
// @Summary Update User
// @Description Update User
// @Tags users
// @ID update-user
// @Accept  json
// @Produce  json
// @Param user body models.UserToUpdate true "Update User Input"
// @Success 200 {object} models.User
// @Router /users [put]
func (ctr *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var userReq models.UserToUpdate
		err := json.Unmarshal(body, &userReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if userReq.Name == "" || userReq.Email == "" || userReq.Password == "" {
			http.Error(w, "Cant be empty values", http.StatusBadRequest)
			return
		}
		userToUpdate, err := models.CreateUser(userReq.Name, userReq.Email, userReq.Password)
		userToUpdate.ID = userReq.ID
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updatedUserID, err := ctr.services.User.UpdateUser(userToUpdate)
		if err != nil || updatedUserID != userToUpdate.ID {
			http.Error(w, "Cant update user", http.StatusInternalServerError)
			return
		}
		updatedUser, _ := ctr.services.User.GetUserByID(updatedUserID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)
	default:
		http.Error(w, "Only PUT method", http.StatusMethodNotAllowed)
	}
}

// DeleteUser godoc
// @Summary Delete User
// @Description Delete User
// @Tags users
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{userId} [delete]
func (ctr *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		id := chi.URLParam(r, "userId")
		user, err := ctr.services.User.GetUserByID(uuid.MustParse(id))
		if err != nil {
			http.Error(w, "User does not exists", http.StatusBadRequest)
			return
		}
		deletedUserID, err := ctr.services.User.DeleteUser(uuid.MustParse(id))
		if err != nil || deletedUserID != uuid.MustParse(id) {
			http.Error(w, "Cant delete user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Only DELETE method", http.StatusMethodNotAllowed)
	}
}
