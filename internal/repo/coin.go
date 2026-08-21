package repo

import (
	"AvitoWinter25/internal/domain"
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

func (r *CoinRepo) GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `
			SELECT 
    				t.id, 
    				uf.username AS user_from, 
    				ut.username AS user_to,
    				t.amount 
			FROM transaction t 
			LEFT JOIN users uf ON uf.id = t.from_user 
			LEFT JOIN users ut ON ut.id = t.to_user
			WHERE from_user = $1 OR to_user = $1`
	rows, err := tr.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction

		if err := rows.Scan(&transaction.ID, &transaction.UserTo, &transaction.Amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
