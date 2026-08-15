package dto

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
