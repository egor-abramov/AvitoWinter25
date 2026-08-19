package domain

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID
	Username       string
	HashedPassword string
	Coins          int
}
