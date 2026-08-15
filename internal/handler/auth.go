package handler

import (
	"AvitoWinter25/internal/handler/dto"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

func NewAuthHandler(authService AuthService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.NewAuthHandler"

		req, ok := dto.Decode[dto.AuthRequest](log, w, r)
		if !ok {
			return
		}

		ctx := r.Context()
		token, err := authService.Login(ctx, req.Username, req.Password)
		if err != nil {
			log.Error(fmt.Sprintf("%s: error processing request: %s", op, err.Error()))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, dto.Error("login user error"))
			return
		}
		resp := dto.AuthResponse{Token: token}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}
