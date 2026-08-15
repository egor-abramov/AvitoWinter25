package main

import (
	"AvitoWinter25/internal/app"
	"AvitoWinter25/internal/config"
	"AvitoWinter25/internal/infrastructure/logger"
	"context"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := logger.Setup()
	cfg := config.MustLoad()

	application, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	application.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Stop(ctx); err != nil {
		log.Error("failed to stop app gracefully", "error", err)
	}

	log.Info("app stopped successfully")
}
