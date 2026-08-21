package merch

import (
	"AvitoWinter25/internal/domain"
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type Repo interface {
	GetMerchByName(ctx context.Context, merchName string) (*domain.Merch, error)
	GetUserMerch(ctx context.Context, userID uuid.UUID) ([]domain.Merch, error)
	AddMerch(ctx context.Context, userID, merchID uuid.UUID) (uuid.UUID, error)
}

type CoinRepo interface {
	AddCoins(ctx context.Context, userID uuid.UUID, amount int) error
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type Service struct {
	repo      Repo
	coinRepo  CoinRepo
	userRepo  UserRepo
	txManager *manager.Manager
}

func New(repo Repo, coinRepo CoinRepo, userRepo UserRepo, txManager *manager.Manager) *Service {
	return &Service{
		repo:      repo,
		coinRepo:  coinRepo,
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (s *Service) GetUserMerch(ctx context.Context, userID uuid.UUID) ([]domain.Merch, error) {
	merch, err := s.repo.GetUserMerch(ctx, userID)
	if err != nil {
		return nil, err
	}
	return merch, nil
}

func (s *Service) Buy(ctx context.Context, userID uuid.UUID, merchName string) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		merch, err := s.repo.GetMerchByName(ctx, merchName)
		if err != nil {
			return err
		}

		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return err
		}
		if merch.Price > user.Coins {
			return fmt.Errorf("merch too expensive")
		}

		if err := s.coinRepo.AddCoins(ctx, user.ID, -merch.Price); err != nil {
			return err
		}

		if _, err := s.repo.AddMerch(ctx, user.ID, merch.ID); err != nil {
			return err
		}
		return nil
	})

	return err
}
