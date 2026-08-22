package info

import (
	"AvitoWinter25/internal/domain"
	"context"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type CoinProvider interface {
	GetCoins(ctx context.Context, userID uuid.UUID) (int, error)
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)
}

type MerchProvider interface {
	GetUserMerch(ctx context.Context, userID uuid.UUID) ([]domain.Merch, error)
}

type Service struct {
	coinProvider  CoinProvider
	merchProvider MerchProvider
}

func New(coinProvider CoinProvider, merchProvider MerchProvider) *Service {
	return &Service{
		coinProvider:  coinProvider,
		merchProvider: merchProvider,
	}
}

func (s *Service) GetInfo(ctx context.Context, userID uuid.UUID) (*domain.Info, error) {
	var (
		coins        int
		transactions []domain.Transaction
		merch        []domain.Merch
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		coins, err = s.coinProvider.GetCoins(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		merch, err = s.merchProvider.GetUserMerch(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		transactions, err = s.coinProvider.GetTransactions(gCtx, userID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &domain.Info{
		Coins:        coins,
		Merch:        merch,
		Transactions: transactions,
	}, nil
}
