package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
)

type TokenController struct {
	services *services.Manager
}

func NewTokenController(services *services.Manager) *TokenController {
	return &TokenController{
		services: services,
	}
}

func (ctr TokenController) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req := &models.UserLoginRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := ctr.services.User.GetUserByEmail(req.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// authenticated

		err = godotenv.Load("config/token.env")
		if err != nil {
			http.Error(w, "Cant load token.env file", http.StatusInternalServerError)
		}
		AccessTokenLifeTime, err := strconv.Atoi(os.Getenv("AccessTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert AccessTokenLifeTime", http.StatusInternalServerError)
		}
		AccessSecret := os.Getenv("AccessSecret")
		RefreshTokenLifeTime, err := strconv.Atoi(os.Getenv("RefreshTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert RefreshTokenLifeTime", http.StatusInternalServerError)
		}
		RefreshSecret := os.Getenv("RefreshSecret")

		accessString, err := token.GenerateToken(user.ID, AccessTokenLifeTime, AccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshString, err := token.GenerateToken(user.ID, RefreshTokenLifeTime, RefreshSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token := &models.Token{
			ID:          uuid.New(),
			UserID:      user.ID,
			AccessHash:  accessString,
			RefreshHash: refreshString,
		}
		_, err = ctr.services.Token.CreateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(token)

	default:
		http.Error(w, "Only Post method", http.StatusMethodNotAllowed)
	}
}

func (ctr *TokenController) Registration(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req := &models.UserRegistrationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := models.CreateUser(req.Name, req.Email, req.Password)
		userID, err := ctr.services.User.CreateUser(user)
		if err != nil {
			http.Error(w, "cant create user", http.StatusInternalServerError)
			return
		}
		// created user

		err = godotenv.Load("config/token.env")
		if err != nil {
			http.Error(w, "Cant load .env file", http.StatusInternalServerError)
		}
		AccessTokenLifeTime, err := strconv.Atoi(os.Getenv("AccessTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert AccessTokenLifeTime", http.StatusInternalServerError)
		}
		AccessSecret := os.Getenv("AccessSecret")
		RefreshTokenLifeTime, err := strconv.Atoi(os.Getenv("RefreshTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert RefreshTokenLifeTime", http.StatusInternalServerError)
		}
		RefreshSecret := os.Getenv("RefreshSecret")

		accessString, err := token.GenerateToken(userID, AccessTokenLifeTime, AccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshString, err := token.GenerateToken(userID, RefreshTokenLifeTime, RefreshSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := &models.Token{
			ID:          uuid.New(),
			UserID:      userID,
			AccessHash:  accessString,
			RefreshHash: refreshString,
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(token)

	default:
		http.Error(w, "Only Post method", http.StatusMethodNotAllowed)
	}

}
