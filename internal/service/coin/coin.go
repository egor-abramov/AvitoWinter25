package coin

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Repo interface {
	AddCoins(ctx context.Context, username string, amount int) error
	GetCoins(ctx context.Context, username string) (int, error)
}

type Service struct {
	log       *slog.Logger
	repo      Repo
	txManager *manager.Manager
}

func New(log *slog.Logger, repo Repo, txManager *manager.Manager) *Service {
	return &Service{
		log:       log,
		repo:      repo,
		txManager: txManager,
	}
}

func (s *Service) GetCoins(ctx context.Context, username string) (int, error) {
	coins, err := s.repo.GetCoins(ctx, username)
	if err != nil {
		return 0, err
	}
	return coins, nil
}

func (s *Service) Transact(ctx context.Context, userFrom string, userTo string, amount int) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		if userTo == userFrom {
			return fmt.Errorf("cannot transact")
		}

		if err := s.repo.AddCoins(ctx, userTo, amount); err != nil {
			return err
		}

		if err := s.repo.AddCoins(ctx, userFrom, -amount); err != nil {
			return err
		}
		return nil
	})

	return err
}
