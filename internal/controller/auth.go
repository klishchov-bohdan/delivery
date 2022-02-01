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

type AuthController struct {
	services *services.Manager
}

func NewAuthController(services *services.Manager) *AuthController {
	return &AuthController{
		services: services,
	}
}

func (ctr AuthController) Login(w http.ResponseWriter, r *http.Request) {
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
		_, _ = ctr.services.Token.DeleteTokenByUserID(user.ID)
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

func (ctr *AuthController) Registration(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req := &models.UserRegistrationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := models.CreateUser(req.Name, req.Email, req.Password)
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
		_, err = ctr.services.User.CreateUserWithTokens(user, newToken)
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

func (ctr *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
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

func (ctr *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := godotenv.Load("config/token.env")
		if err != nil {
			http.Error(w, "Cant load .env file", http.StatusInternalServerError)
		}
		accessTokenLifeTime, err := strconv.Atoi(os.Getenv("AccessTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert AccessTokenLifeTime", http.StatusInternalServerError)
		}
		accessSecret := os.Getenv("AccessSecret")
		refreshTokenLifeTime, err := strconv.Atoi(os.Getenv("RefreshTokenLifeTime"))
		if err != nil {
			http.Error(w, "Cant convert RefreshTokenLifeTime", http.StatusInternalServerError)
		}
		refreshSecret := os.Getenv("RefreshSecret")
		refreshString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		isValid, err := token.ValidateToken(refreshString, refreshSecret)
		if err != nil || !isValid {
			http.Error(w, "refresh: invalid token", http.StatusUnauthorized)
			return
		}
		claims, err := token.GetClaims(refreshString, refreshSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenDB, err := ctr.services.Token.GetTokenByUserID(claims.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tokenDB.RefreshHash != base64.StdEncoding.EncodeToString([]byte(refreshString)) {
			http.Error(w, "refresh: invalid token", http.StatusInternalServerError)
			return
		}
		newAccessString, err := token.GenerateToken(claims.ID, accessTokenLifeTime, accessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newRefreshString, err := token.GenerateToken(claims.ID, refreshTokenLifeTime, refreshSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := &models.ResponseToken{
			Access:  newAccessString,
			Refresh: newRefreshString,
		}

		tokenDB.AccessHash = base64.StdEncoding.EncodeToString([]byte(newAccessString))
		tokenDB.RefreshHash = base64.StdEncoding.EncodeToString([]byte(newRefreshString))
		_, err = ctr.services.Token.UpdateToken(tokenDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}

}
