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
			println("here 1")
			http.Error(w, "Cant load token.env file", http.StatusInternalServerError)
			return
		}
		accessSecret := os.Getenv("AccessSecret")
		claims, err := token.GetClaims(accessString, accessSecret)
		if err != nil {
			println("here 2")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		isValid, err := token.ValidateToken(accessString, accessSecret)
		if err != nil {
			println("here 3")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isValid {
			println("here 4")
			_, err := m.service.Token.DeleteTokenByUserID(claims.ID)
			if err != nil {
				println("here 5")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "middleware: invalid token", http.StatusUnauthorized)
			return
		}
		received, err := m.service.Token.GetTokenByUserID(claims.ID)
		if err != nil {
			println("here 6")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if received == nil {
			println("here 7")
			http.Error(w, "token does not exists", http.StatusUnauthorized)
			return
		}
		if base64.StdEncoding.EncodeToString([]byte(accessString)) != received.AccessHash {
			println("here 8")
			http.Error(w, "middleware: invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
