package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/ParkPawapon/request-api/internal/app"
	"github.com/ParkPawapon/request-api/internal/config"
	applogger "github.com/ParkPawapon/request-api/internal/infrastructure/logger"
	"go.uber.org/zap"
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

	logger.Info("runtime configuration loaded",
		zap.String("app_name", cfg.App.Name),
		zap.String("app_env", cfg.App.Env),
		zap.String("app_port", cfg.App.Port),
		zap.Duration("request_timeout", cfg.App.RequestTimeout),
		zap.Int64("max_body_bytes", cfg.App.MaxBodyBytes),
		zap.Bool("database_url_configured", cfg.Database.URL != ""),
		zap.Int("database_max_open_conns", cfg.Database.MaxOpenConns),
		zap.Int("database_max_idle_conns", cfg.Database.MaxIdleConns),
		zap.String("redis_addr", cfg.Redis.Addr),
		zap.Int("redis_db", cfg.Redis.DB),
		zap.String("redis_key_prefix", cfg.Redis.KeyPrefix),
		zap.Int("cors_allowed_origins", len(cfg.CORS.AllowedOrigins)),
	)

	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("application startup failed", applogger.Error(err))
	}
	defer application.Close(ctx)

	if err := application.Run(ctx); err != nil {
		logger.Fatal("application stopped with error", applogger.Error(err))
	}
}
