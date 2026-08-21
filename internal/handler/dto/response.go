package dto

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedCoins struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoins struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedCoins `json:"received"`
	Sent     []SentCoins     `json:"sent"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Errors: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) ErrorResponse {
	var errMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is required", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is not valid", err.Field()))
		}
	}
	return Error(strings.Join(errMessages, "; "))
}

func RespondWithError(w http.ResponseWriter, r *http.Request, log *slog.Logger, op, msg string, err error, status int) {
	log.Error("request failed", slog.String("op", op), slog.String("msg", msg), slog.Any("err", err))

	render.Status(r, status)
	render.JSON(w, r, Error(msg))
}
