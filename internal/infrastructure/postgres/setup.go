package postgres

import (
	"AvitoWinter25/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(cfg config.DB) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgxConfig.MaxConns = 20
	pgxConfig.MinConns = 2
	pgxConfig.MaxConnLifetime = 1 * time.Hour
	pgxConfig.MaxConnIdleTime = 15 * time.Minute
	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
