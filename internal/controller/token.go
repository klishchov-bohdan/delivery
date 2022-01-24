package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
		_, err = ctr.services.Token.DeleteTokenByUserID(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// authenticated
		err = godotenv.Load("config/token.env")
		if err != nil {
			http.Error(w, "Cant load token.env file", http.StatusInternalServerError)
			return
		}
		AccessTokenLifeTime, err := strconv.Atoi(os.Getenv("AccessTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert AccessTokenLifeTime", http.StatusInternalServerError)
			return
		}
		AccessSecret := os.Getenv("AccessSecret")
		RefreshTokenLifeTime, err := strconv.Atoi(os.Getenv("RefreshTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert RefreshTokenLifeTime", http.StatusInternalServerError)
			return
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
		resp := &models.ResponseToken{
			Access:  accessString,
			Refresh: refreshString,
		}
		newToken := &models.Token{
			ID:          uuid.New(),
			UserID:      user.ID,
			AccessHash:  base64.StdEncoding.EncodeToString([]byte(accessString)),
			RefreshHash: base64.StdEncoding.EncodeToString([]byte(refreshString)),
		}
		_, err = ctr.services.Token.CreateToken(newToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		resp := &models.ResponseToken{
			Access:  accessString,
			Refresh: refreshString,
		}

		newToken := &models.Token{
			ID:          uuid.New(),
			UserID:      user.ID,
			AccessHash:  base64.StdEncoding.EncodeToString([]byte(accessString)),
			RefreshHash: base64.StdEncoding.EncodeToString([]byte(refreshString)),
		}
		_, err = ctr.services.Token.CreateToken(newToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Only Post method", http.StatusMethodNotAllowed)
	}

}

func (ctr *TokenController) Profile(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

func (ctr *TokenController) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		accessString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		accessSecret := os.Getenv("AccessSecret")
		claims, err := token.GetClaims(accessString, accessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		uid, err := ctr.services.Token.DeleteTokenByUserID(claims.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Logout user %s", uid.String())))
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}
