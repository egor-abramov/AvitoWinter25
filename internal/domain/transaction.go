package domain

import "github.com/google/uuid"

type Transaction struct {
	ID       uuid.UUID
	UserFrom User
	UserTo   User
	Amount   int
}
