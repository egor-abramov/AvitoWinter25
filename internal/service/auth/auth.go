package auth

import (
	"AvitoWinter25/internal/domain"
	"context"
	"errors"
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

func New(repo Repo, log *slog.Logger, jwtSecret string, jwrExpire time.Duration) *Service {
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
		s.log.Error(fmt.Sprintf("%s: error getting user by username: %s", op, err.Error()))
		return "", err
	}

	if user == nil {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			s.log.Error(fmt.Sprintf("%s: error hashing password: %s", op, err.Error()))
			return "", err
		}
		id, err := s.repo.CreateUser(ctx, username, hashedPassword)
		if err != nil {
			s.log.Error(fmt.Sprintf("%s: error creating user: %s", op, err.Error()))
			return "", err
		}
		return generateToken(id, username, s.jwtSecret, s.jwtExpire)
	}

	if !checkPasswordHash(password, user.HashedPassword) {
		s.log.Error(fmt.Sprintf("%s: error hashing password: %s", op, err.Error()))
		return "", errors.New("invalid password")
	}
	return generateToken(user.ID, username, s.jwtSecret, s.jwtExpire)
}
