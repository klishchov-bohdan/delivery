package middleware

import (
	"encoding/base64"
	"github.com/joho/godotenv"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"net/http"
	"os"
)

type Middleware struct {
	service *services.Manager
}

func NewMiddleware(service *services.Manager) *Middleware {
	return &Middleware{
		service: service,
	}
}

func (m *Middleware) AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessString := token.GetTokenFromBearerString(r.Header.Get("Authorization"))
		err := godotenv.Load("config/token.env")
		if err != nil {
			http.Error(w, "Cant load token.env file", http.StatusInternalServerError)
			return
		}
		accessSecret := os.Getenv("AccessSecret")
		isValid, err := token.ValidateToken(accessString, accessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isValid {
			http.Error(w, "middleware: invalid token", http.StatusUnauthorized)
			return
		}
		claims, err := token.GetClaims(accessString, accessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		received, err := m.service.Token.GetTokenByUserID(claims.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if received == nil {
			http.Error(w, "user has not authorized", http.StatusUnauthorized)
			return
		}
		if base64.StdEncoding.EncodeToString([]byte(accessString)) != received.AccessHash {
			http.Error(w, "middleware: user has not authorised", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
