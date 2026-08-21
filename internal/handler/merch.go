package handler

import (
	"AvitoWinter25/internal/handler/dto"
	"AvitoWinter25/internal/handler/middleware"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type BuyService interface {
	Buy(ctx context.Context, userID uuid.UUID, merchName string) error
}

func NewBuyHandler(service BuyService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.NewBuyHandler"

		merch := chi.URLParam(r, "item")

		ctx := r.Context()
		userID, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			msg := "user not found in context"
			dto.RespondWithError(w, r, log, op, msg, nil, http.StatusInternalServerError)
			return
		}

		err := service.Buy(ctx, userID, merch)
		if err != nil {
			msg := "buying error"
			dto.RespondWithError(w, r, log, op, msg, err, http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
