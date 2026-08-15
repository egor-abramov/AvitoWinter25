package domain

import "github.com/google/uuid"

type Merch struct {
	ID    uuid.UUID
	Name  string
	Price int
}
