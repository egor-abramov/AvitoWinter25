package handler

import (
	"AvitoWinter25/internal/domain"
	"AvitoWinter25/internal/handler/dto"
	"AvitoWinter25/internal/handler/middleware"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type TransactionService interface {
	Transact(ctx context.Context, t domain.Transaction) error
}

func NewSendCoinHandler(service TransactionService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.NewSendCoinHandler"

		req, ok := dto.Decode[dto.SendCoinRequest](log, w, r)
		if !ok {
			return
		}

		ctx := r.Context()
		userFrom, ok := ctx.Value(middleware.UserName).(string)
		if !ok {
			msg := "username not found in context"
			dto.RespondWithError(w, r, log, op, msg, http.StatusInternalServerError)
			return
		}

		tr := domain.Transaction{
			UserFrom: userFrom,
			UserTo:   req.UserTo,
			Amount:   req.Amount,
		}
		err := service.Transact(ctx, tr)
		if err != nil {
			dto.RespondWithError(w, r, log, op, err.Error(), http.StatusBadRequest)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
