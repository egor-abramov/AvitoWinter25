package app

import (
	"AvitoWinter25/internal/config"
	"AvitoWinter25/internal/handler"
	"AvitoWinter25/internal/handler/middleware"
	"AvitoWinter25/internal/infrastructure/postgres"
	"AvitoWinter25/internal/repo"
	"AvitoWinter25/internal/service/auth"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
	log    *slog.Logger
	pool   *pgxpool.Pool
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	pool, err := postgres.Setup(cfg.DB)
	if err != nil {
		return nil, err
	}

	userRepo := repo.NewUserRepo(pool, nil)
	authService := auth.New(userRepo, log, cfg.JWT.Secret, cfg.JWT.ExpTime)
	authHandler := handler.NewAuthHandler(authService, log)

	r := chi.NewRouter()
	r.Post("/api/auth", authHandler)
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewJWTExtractor(cfg.JWT.Secret))
	})

	srv := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: r,
	}

	return &App{
		server: srv,
		log:    log,
		pool:   pool,
	}, nil
}

func (a *App) Run() {
	a.log.Info("starting http server", "address", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("failed to start server", "error", err)
		}
	}()
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("stopping http server")

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.log.Info("closing database connection pool")
	a.pool.Close()

	return nil
}
