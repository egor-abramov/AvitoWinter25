package repo

import (
	"AvitoWinter25/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, hashed_password, coins from users WHERE username = $1`
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Balance,
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
	var id uuid.UUID

	query := `INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id`
	err := r.pool.QueryRow(ctx, query, username, hashedPassword).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
