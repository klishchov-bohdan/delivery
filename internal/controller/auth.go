package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	services *services.Manager
}

func NewAuthController(services *services.Manager) *AuthController {
	return &AuthController{
		services: services,
	}
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags auth
// @ID user-login
// @Accept  json
// @Produce  json
// @Param authLogin body models.UserLoginRequest true "Auth Login Input"
// @Success 200 {object} models.ResponseToken
// @Router /login [post]
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
		cfg := config.NewConfig()

		accessString, err := token.GenerateToken(user.ID, cfg.AccessTokenLifeTime, cfg.AccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshString, err := token.GenerateToken(user.ID, cfg.RefreshTokenLifeTime, cfg.RefreshSecret)
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

// Registration godoc
// @Summary Registration
// @Description Registration
// @Tags auth
// @ID user-registration
// @Accept  json
// @Produce  json
// @Param authRegistration body models.UserRegistrationRequest true "Auth Registration Input"
// @Success 200 {object} models.ResponseToken
// @Router /registration [post]
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
		cfg := config.NewConfig()
		accessString, err := token.GenerateToken(user.ID, cfg.AccessTokenLifeTime, cfg.AccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshString, err := token.GenerateToken(user.ID, cfg.RefreshTokenLifeTime, cfg.RefreshSecret)
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

// Logout godoc
// @Summary Logout
// @Description Logout
// @Tags auth
// @ID user-logout
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization"
// @Success 200
// @Security ApiKeyAuth
// @Router /logout [post]
func (ctr *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		accessString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		cfg := config.NewConfig()
		claims, err := token.GetClaims(accessString, cfg.AccessSecret)
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

// Refresh godoc
// @Summary Refresh
// @Description Refresh
// @Tags auth
// @ID user-refresh
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization"
// @Success 200 {object} models.ResponseToken
// @Security ApiKeyAuth
// @Router /refresh [post]
func (ctr *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		cfg := config.NewConfig()
		refreshString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		isValid, err := token.ValidateToken(refreshString, cfg.RefreshSecret)
		if err != nil || !isValid {
			http.Error(w, "refresh: invalid token", http.StatusUnauthorized)
			return
		}
		claims, err := token.GetClaims(refreshString, cfg.RefreshSecret)
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
		newAccessString, err := token.GenerateToken(claims.ID, cfg.AccessTokenLifeTime, cfg.AccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newRefreshString, err := token.GenerateToken(claims.ID, cfg.RefreshTokenLifeTime, cfg.RefreshSecret)
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
