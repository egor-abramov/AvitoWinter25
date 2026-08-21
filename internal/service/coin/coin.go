package coin

import (
	"AvitoWinter25/internal/domain"
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type Repo interface {
	AddCoins(ctx context.Context, username string, amount int) error
	GetCoins(ctx context.Context, username string) (int, error)
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)
}

type Service struct {
	repo      Repo
	txManager *manager.Manager
}

func New(repo Repo, txManager *manager.Manager) *Service {
	return &Service{
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

func (s *Service) Transact(ctx context.Context, t domain.Transaction) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		if t.UserFrom == t.UserTo {
			return fmt.Errorf("cannot transact")
		}

		if err := s.repo.AddCoins(ctx, t.UserFrom, t.Amount); err != nil {
			return err
		}

		if err := s.repo.AddCoins(ctx, t.UserTo, -t.Amount); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	transactions, err := s.repo.GetTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
