package middleware

import (
	"encoding/base64"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"net/http"
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
		if accessString == "" {
			http.Error(w, "middleware: invalid bearer string", http.StatusUnauthorized)
			return
		}
		cfg := config.NewConfig()
		isValid, err := token.ValidateToken(accessString, cfg.AccessSecret)
		if err != nil {
			http.Error(w, "middleware: invalid token", http.StatusUnauthorized)
			return
		}
		if !isValid {
			http.Error(w, "middleware: invalid token", http.StatusUnauthorized)
			return
		}
		claims, err := token.GetClaims(accessString, cfg.AccessSecret)
		if err != nil {
			http.Error(w, "middleware: can`t cet claims", http.StatusUnauthorized)
			return
		}
		received, err := m.service.Token.GetTokenByUserID(claims.ID)
		if received == nil || err != nil {
			http.Error(w, "middleware: user has not authorized", http.StatusUnauthorized)
			return
		}
		if base64.StdEncoding.EncodeToString([]byte(accessString)) != received.AccessHash {
			http.Error(w, "middleware: user has not authorised", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
