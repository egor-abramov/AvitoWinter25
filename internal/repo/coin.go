package repo

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type CoinRepo struct {
	conn pgx.Conn
	log  *slog.Logger
}

func NewCoinRepo(conn pgx.Conn, log *slog.Logger) *CoinRepo {
	return &CoinRepo{
		conn: conn,
		log:  log,
	}
}
