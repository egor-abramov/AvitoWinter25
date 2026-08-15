package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateToken(id uuid.UUID, username, secret string, expire time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(expire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
