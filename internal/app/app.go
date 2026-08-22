package app

import (
	"AvitoWinter25/internal/config"
	"AvitoWinter25/internal/handler"
	"AvitoWinter25/internal/handler/middleware"
	"AvitoWinter25/internal/infrastructure/postgres"
	"AvitoWinter25/internal/repo"
	"AvitoWinter25/internal/service/auth"
	"AvitoWinter25/internal/service/coin"
	"AvitoWinter25/internal/service/info"
	"AvitoWinter25/internal/service/merch"
	"context"
	"errors"
	"log/slog"
	"net/http"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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

	userRepo := repo.NewUserRepo(pool)
	coinRepo := repo.NewCoinRepo(pool)
	merchRepo := repo.NewMerchRepo(pool)

	trFactory := trmpgx.NewFactory(pool)
	txManager, err := manager.New(trFactory)
	if err != nil {
		return nil, err
	}

	authService := auth.New(userRepo, log, cfg.JWT.Secret, cfg.JWT.ExpTime)
	coinService := coin.New(coinRepo, txManager)
	merchService := merch.New(merchRepo, coinRepo, userRepo, txManager)
	infoService := info.New(coinService, merchService)

	authHandler := handler.NewAuthHandler(authService, log)
	coinHandler := handler.NewSendCoinHandler(coinService, log)
	merchHandler := handler.NewBuyHandler(merchService, log)
	infoHandler := handler.NewInfoHandler(infoService, log)

	r := chi.NewRouter()
	r.Post("/api/auth", authHandler)
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewJWTExtractor(cfg.JWT.Secret))
		r.Get("/api/info", infoHandler)
		r.Post("/api/sendCoin", coinHandler)
		r.Get("/api/buy/{item}", merchHandler)
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
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
