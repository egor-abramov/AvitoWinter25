package repo

import (
	"AvitoWinter25/internal/domain"
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewMerchRepo(pool *pgxpool.Pool) *MerchRepo {
	return &MerchRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func (r *MerchRepo) GetMerchByName(ctx context.Context, merchName string) (*domain.Merch, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var merch domain.Merch

	query := `SELECT id, name, price from merch WHERE name = $1`
	err := tr.QueryRow(ctx, query, merchName).Scan(
		&merch.ID,
		&merch.Name,
		&merch.Price,
	)
	if err != nil {
		return nil, err
	}

	return &merch, nil
}

func (r *MerchRepo) GetUserMerch(ctx context.Context, userID uuid.UUID) ([]domain.Merch, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `SELECT m.id, m.name, um.quantity FROM merch m JOIN user_merch um ON m.id = um.merch_id WHERE um.user_id = $1`

	rows, err := tr.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchItems []domain.Merch
	for rows.Next() {
		var merch domain.Merch
		if err := rows.Scan(&merch.ID, &merch.Name, &merch.Quantity); err != nil {
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

	query := `
				INSERT INTO user_merch 
				    (user_id, merch_id, quantity) 
				VALUES ($1, $2, 1) 
				ON CONFLICT (user_id, merch_id)
				DO UPDATE SET quantity = user_merch.quantity + EXCLUDED.quantity 
				RETURNING id`
	err := tr.QueryRow(ctx, query, userID, merchID).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
