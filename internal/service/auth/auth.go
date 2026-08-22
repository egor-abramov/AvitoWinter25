package auth

import (
	"AvitoWinter25/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateUser(ctx context.Context, username, hashedPassword string) (uuid.UUID, error)
}

type Service struct {
	log       *slog.Logger
	repo      Repo
	jwtSecret string
	jwtExpire time.Duration
}

func New(repo Repo, jwtSecret string, jwrExpire time.Duration, log *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		log:       log,
		jwtSecret: jwtSecret,
		jwtExpire: jwrExpire,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	const op = "service.auth.Login"

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.Error("error getting user by username", slog.String("op", op), slog.Any("err", err))
		return "", fmt.Errorf("error getting user by username: %s", username)
	}

	if user == nil {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			s.log.Error("error hashing password", slog.String("op", op), slog.Any("err", err))
			return "", fmt.Errorf("error hashing password")
		}
		id, err := s.repo.CreateUser(ctx, username, hashedPassword)
		if err != nil {
			s.log.Error("error creating user", slog.String("op", op), slog.Any("err", err))
			return "", fmt.Errorf("error creating user")
		}
		return generateToken(id, username, s.jwtSecret, s.jwtExpire)
	}

	if !checkPasswordHash(password, user.HashedPassword) {
		s.log.Error("invalid password", slog.String("op", op), slog.Any("err", err))
		return "", fmt.Errorf("invalid password")
	}
	return generateToken(user.ID, username, s.jwtSecret, s.jwtExpire)
}
