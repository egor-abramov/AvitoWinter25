package repo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinRepo struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewCoinRepo(pool *pgxpool.Pool, log *slog.Logger) *CoinRepo {
	return &CoinRepo{
		pool: pool,
		log:  log,
	}
}

func (r *CoinRepo) AddCoins(ctx context.Context, username string, amount int) error {
	query := `UPDATE users SET coins = coins + $1 WHERE username = $2`

	_, err := r.pool.Exec(ctx, query, amount, username)
	return err
}

func (r *CoinRepo) GetCoins(ctx context.Context, username string) (int, error) {
	var coins int

	query := `SELECT coins from users WHERE username = $1`
	err := r.pool.QueryRow(ctx, query, username).Scan(&coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return coins, nil
}
