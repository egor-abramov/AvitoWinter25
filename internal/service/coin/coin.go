package coin

import (
	"AvitoWinter25/internal/domain"
	"context"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type Repo interface {
	AddCoinsByUsername(ctx context.Context, username string, amount int) error
	GetCoins(ctx context.Context, userID uuid.UUID) (int, error)
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)
	AddTransaction(ctx context.Context, userFrom *domain.User, userTo *domain.User, amount int) (uuid.UUID, error)
}

type UserRepo interface {
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}

type Service struct {
	repo      Repo
	userRepo  UserRepo
	txManager *manager.Manager
	log       *slog.Logger
}

func New(repo Repo, userRepo UserRepo, txManager *manager.Manager, log *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		userRepo:  userRepo,
		txManager: txManager,
		log:       log,
	}
}

func (s *Service) GetCoins(ctx context.Context, userID uuid.UUID) (int, error) {
	const op = "service.coin.GetCoins"

	coins, err := s.repo.GetCoins(ctx, userID)
	if err != nil {
		s.log.Error("getting coins failed", slog.String("op", op), slog.Any("err", err))
		return 0, fmt.Errorf("cannot get coins")
	}
	return coins, nil
}

func (s *Service) Transact(ctx context.Context, t domain.Transaction) error {
	const op = "service.coin.Transact"

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		isUserExists, err := s.userRepo.IsUsernameExists(ctx, t.UserTo)
		if err != nil {
			return err
		}
		if !isUserExists {
			s.log.Error(fmt.Sprintf("user %s does not exist", t.UserTo), slog.String("op", op))
			return fmt.Errorf("user %s does not exist", t.UserTo)
		}

		if t.Amount <= 0 {
			s.log.Error("amount must be greater than zero", slog.String("op", op))
			return fmt.Errorf("amount must be greater than zero")
		}

		if t.UserFrom == t.UserTo {
			s.log.Error("userTo can't be userFrom", slog.String("op", op))
			return fmt.Errorf("cannot transact")
		}

		if err := s.repo.AddCoinsByUsername(ctx, t.UserFrom, -t.Amount); err != nil {
			s.log.Error("adding coins failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("cannot add coins by username")
		}

		if err := s.repo.AddCoinsByUsername(ctx, t.UserTo, t.Amount); err != nil {
			s.log.Error("adding coins failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("cannot add coins by username")
		}

		userFrom, err := s.userRepo.GetUserByUsername(ctx, t.UserFrom)
		if err != nil {
			s.log.Error("getting userFrom failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("cannot get user from")
		}

		userTo, err := s.userRepo.GetUserByUsername(ctx, t.UserTo)
		if err != nil {
			s.log.Error("getting userTo failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("cannot get user to")
		}

		if _, err := s.repo.AddTransaction(ctx, userFrom, userTo, t.Amount); err != nil {
			s.log.Error("add transaction failed", slog.String("op", op), slog.Any("err", err))
			return fmt.Errorf("cannot add transaction to history")
		}
		return nil
	})

	return err
}

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	const op = "service.coin.GetTransactions"

	transactions, err := s.repo.GetTransactions(ctx, userID)
	if err != nil {
		s.log.Error("getting transaction failed", slog.String("op", op), slog.Any("err", err))
		return nil, fmt.Errorf("cannot get transactions: %w", err)
	}
	return transactions, nil
}
