package repo

import (
	"AvitoWinter25/internal/domain"
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewCoinRepo(pool *pgxpool.Pool) *CoinRepo {
	return &CoinRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func (r *CoinRepo) AddCoins(ctx context.Context, userID uuid.UUID, amount int) error {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `UPDATE users SET coins = coins + $1 WHERE id = $2`

	_, err := tr.Exec(ctx, query, amount, userID)
	return err
}

func (r *CoinRepo) AddCoinsByUsername(ctx context.Context, username string, amount int) error {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `UPDATE users SET coins = coins + $1 WHERE username = $2`

	_, err := tr.Exec(ctx, query, amount, username)
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

func (r *CoinRepo) AddTransaction(ctx context.Context, userFrom *domain.User, userTo *domain.User, amount int) (uuid.UUID, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `INSERT INTO transaction (from_user, to_user, amount) VALUES ($1, $2, $3) RETURNING id`
	var id uuid.UUID
	err := tr.QueryRow(ctx, query, userFrom.ID, userTo.ID, amount).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
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

		if err := rows.Scan(&transaction.ID, &transaction.UserFrom, &transaction.UserTo, &transaction.Amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
