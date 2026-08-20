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

type MerchRepo struct {
	pool   *pgxpool.Pool
	log    *slog.Logger
	getter *trmpgx.CtxGetter
}

func NewMerchRepo(log *slog.Logger, pool *pgxpool.Pool, getter *trmpgx.CtxGetter) *MerchRepo {
	return &MerchRepo{
		pool:   pool,
		log:    log,
		getter: getter,
	}
}

func (r *MerchRepo) GetMerchByName(ctx context.Context, merchName string) (*domain.Merch, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var merch domain.Merch

	query := `SELECT id, name, price from merch WHERE merch = $1`
	err := tr.QueryRow(ctx, query, merchName).Scan(
		&merch.ID,
		&merch.Name,
		&merch.Price,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return &merch, err
	}

	return &merch, nil
}

func (r *MerchRepo) GetUserMerch(ctx context.Context, userID uuid.UUID) ([]domain.Merch, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `SELECT m.id, m.name FROM merch m JOIN user_merch um ON m.id = um.merch_id WHERE um.user_id = $1`

	rows, err := tr.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchItems []domain.Merch
	for rows.Next() {
		var merch domain.Merch
		if err := rows.Scan(&merch.ID, &merch.Name); err != nil {
			return nil, err
		}
		merchItems = append(merchItems, merch)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return merchItems, nil
}

func (r *MerchRepo) AddMerch(ctx context.Context, userID, merchID uuid.UUID) (uuid.UUID, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var id uuid.UUID

	query := `INSERT INTO user_merch (user_id, merch_id) VALUES ($1, $2) RETURNING id`
	err := tr.QueryRow(ctx, query, userID, merchID).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
