package domain

import "github.com/google/uuid"

type Transaction struct {
	ID       uuid.UUID
	UserFrom string
	UserTo   string
	Amount   int
}
