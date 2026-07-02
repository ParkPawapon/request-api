package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/ParkPawapon/request-api/internal/app"
	"github.com/ParkPawapon/request-api/internal/config"
	applogger "github.com/ParkPawapon/request-api/internal/infrastructure/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	logger, err := applogger.New(cfg.App.Env, cfg.Log.Level)
	if err != nil {
		log.Fatalf("logger error: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("application startup failed", applogger.Error(err))
	}
	defer application.Close(ctx)

	if err := application.Run(ctx); err != nil {
		logger.Fatal("application stopped with error", applogger.Error(err))
	}
}
