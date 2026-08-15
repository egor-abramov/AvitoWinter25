package dto

type inventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type receivedCoins struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type sentCoins struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []receivedCoins `json:"received"`
	Sent     []sentCoins     `json:"sent"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []inventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
