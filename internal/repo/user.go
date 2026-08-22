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

type UserRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func (r *UserRepo) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	query := `SELECT EXISTS (SELECT 1 from users WHERE username = $1)`
	var exists bool
	err := tr.QueryRow(ctx, query, username).Scan(&exists)
	return exists, err
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var user domain.User

	query := `SELECT id, username, hashed_password, coins from users WHERE username = $1`
	err := tr.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Coins,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var user domain.User

	query := `SELECT id, username, hashed_password, coins from users WHERE id = $1`
	err := tr.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Coins,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, username, hashedPassword string) (uuid.UUID, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.pool)

	var id uuid.UUID

	query := `INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id`
	err := tr.QueryRow(ctx, query, username, hashedPassword).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
