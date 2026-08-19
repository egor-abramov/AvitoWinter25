package coin

import (
	"context"
	"log/slog"
)

type Repo interface {
	AddCoins(ctx context.Context, username string, amount int) error
	GetCoins(ctx context.Context, username string) (int, error)
}

type Service struct {
	log  *slog.Logger
	repo Repo
}

func New(log *slog.Logger, repo Repo) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) GetCoins(ctx context.Context, username string) (int, error) {
	panic("implement me")
}

func (s *Service) BuyMerch(ctx context.Context, username, merch string) error {
	panic("implement me")
}

func (s *Service) Transact(ctx context.Context, userFrom string, userTo string, amount int) error {
	panic("implement me")
}
