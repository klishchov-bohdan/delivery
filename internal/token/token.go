package token

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type JWTCustomClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(userID uuid.UUID, lifetimeMinutes int, secret string) (string, error) {
	claims := &JWTCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(lifetimeMinutes)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GetTokenFromBearerString(bearer string) string {
	if bearer == "" {
		return ""
	}
	parts := strings.Split(bearer, "Bearer ")
	if len(parts) != 2 {
		return ""
	}
	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}
	return token
}

func GetClaims(tokenStr string, secret string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTCustomClaims{}, func(token *jwt.Token) (
		interface{}, error,
	) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}
	return claims, nil
}

func ValidateToken(tokenStr string, secret string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (
		interface{}, error,
	) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, errors.New("invalid token")
	}
	return true, nil
}
