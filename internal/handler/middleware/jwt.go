package middleware

import (
	"AvitoWinter25/internal/handler/dto"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey string = "user_id"

func NewJWTExtractor(secretKey string) func(handler http.Handler) http.Handler {
	const op = "handler.middleware.NewJWTExtractor"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, dto.Error("missing authorization header"))
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, dto.Error("invalid authorization header format"))
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
				render.JSON(w, r, dto.Error("invalid or expired token"))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, dto.Error("failed to parse token claims"))
				return
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, dto.Error("user_id not found in token"))
				return
			}
			userID := int(userIDFloat)

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
