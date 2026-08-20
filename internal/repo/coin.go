package repo

import (
	"context"
	"errors"
	"log/slog"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinRepo struct {
	pool   *pgxpool.Pool
	log    *slog.Logger
	getter *trmpgx.CtxGetter
}

func NewCoinRepo(log *slog.Logger, pool *pgxpool.Pool, getter *trmpgx.CtxGetter) *CoinRepo {
	return &CoinRepo{
		pool:   pool,
		log:    log,
		getter: getter,
	}
}

func (r *CoinRepo) AddCoins(ctx context.Context, userID uuid.UUID, amount int) error {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `UPDATE users SET coins = coins + $1 WHERE id = $2`

	_, err := tr.Exec(ctx, query, amount, userID)
	return err
}

func (r *CoinRepo) GetCoins(ctx context.Context, userID uuid.UUID) (int, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var coins int

	query := `SELECT coins from users WHERE id = $1`
	err := tr.QueryRow(ctx, query, userID).Scan(&coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return coins, nil
}
