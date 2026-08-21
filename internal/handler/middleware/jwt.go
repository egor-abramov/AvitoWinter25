package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const UserIDKey string = "id"
const UserName string = "username"

func NewJWTExtractor(secretKey string) func(handler http.Handler) http.Handler {
	const op = "handler.middleware.NewJWTExtractor"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				return
			}
			tokenStr := headerParts[1]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				render.Status(r, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(claims["id"].(string))
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				return
			}

			username, ok := claims["username"].(string)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserName, username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
