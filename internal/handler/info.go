package handler

import (
	"AvitoWinter25/internal/domain"
	"AvitoWinter25/internal/handler/dto"
	"AvitoWinter25/internal/handler/middleware"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type InfoService interface {
	GetInfo(ctx context.Context, userID uuid.UUID, username string) (*domain.Info, error)
}

func NewInfoHandler(infoService InfoService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.NewInfoHandler"

		ctx := r.Context()
		id, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			msg := "user id not found in context"
			dto.RespondWithError(w, r, log, op, msg, nil, http.StatusInternalServerError)
			return
		}

		username, ok := ctx.Value(middleware.UserName).(string)
		if !ok {
			msg := "username not found in context"
			dto.RespondWithError(w, r, log, op, msg, nil, http.StatusInternalServerError)
			return
		}

		info, err := infoService.GetInfo(ctx, id, username)
		if err != nil {
			msg := "error getting info"
			dto.RespondWithError(w, r, log, op, msg, err, http.StatusInternalServerError)
			return
		}

		infoResp := dto.InfoResponse{
			Coins: info.Coins,
		}

		inventory := make([]dto.InventoryItem, 0, len(info.Merch))
		for _, m := range info.Merch {
			item := dto.InventoryItem{
				Quantity: m.Quantity,
				Type:     m.Name,
			}
			inventory = append(inventory, item)
		}
		infoResp.Inventory = inventory

		var history dto.CoinHistory
		history.Sent = make([]dto.SentCoins, 0, len(info.Transactions))
		history.Received = make([]dto.ReceivedCoins, 0, len(info.Transactions))
		for _, h := range info.Transactions {
			if h.UserFrom == username {
				item := dto.SentCoins{
					ToUser: h.UserTo,
					Amount: h.Amount,
				}
				history.Sent = append(history.Sent, item)
			} else {
				item := dto.ReceivedCoins{
					FromUser: h.UserFrom,
					Amount:   h.Amount,
				}
				history.Received = append(history.Received, item)
			}
		}
		infoResp.CoinHistory = history

		render.Status(r, http.StatusOK)
		render.JSON(w, r, infoResp)
	}
}
