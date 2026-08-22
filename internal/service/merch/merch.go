package merch

import (
	"AvitoWinter25/internal/domain"
	"context"
	"fmt"
	"log/slog"

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
	log       *slog.Logger
}

func New(repo Repo, coinRepo CoinRepo, userRepo UserRepo, txManager *manager.Manager, log *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		coinRepo:  coinRepo,
		userRepo:  userRepo,
		txManager: txManager,
		log:       log,
	}
}

func (s *Service) GetUserMerch(ctx context.Context, userID uuid.UUID) ([]domain.Merch, error) {
	const op = "service.merch.GetUserMerch"

	merch, err := s.repo.GetUserMerch(ctx, userID)
	if err != nil {
		s.log.Error("getting merch failed", slog.String("op", op), slog.Any("err", err))
		return nil, fmt.Errorf("getting merch failed")
	}
	return merch, nil
}

func (s *Service) Buy(ctx context.Context, userID uuid.UUID, merchName string) error {
	const op = "service.merch.Buy"

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		merch, err := s.repo.GetMerchByName(ctx, merchName)
		if err != nil {
			s.log.Error("getting merch by name failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("getting merch by name failed")
		}

		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			s.log.Error("buying merch failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("buying merch failed")
		}
		if merch.Price > user.Coins {
			s.log.Error("buying merch failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("merch too expensive")
		}

		if err := s.coinRepo.AddCoins(ctx, user.ID, -merch.Price); err != nil {
			s.log.Error("adding merch failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("adding merch failed")
		}

		if _, err := s.repo.AddMerch(ctx, user.ID, merch.ID); err != nil {
			s.log.Error("adding merch failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("adding merch failed")
		}
		return nil
	})

	return err
}
