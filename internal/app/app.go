package app

import (
	"context"
	"net/http"

	"github.com/ParkPawapon/request-api/internal/config"
	"github.com/ParkPawapon/request-api/internal/infrastructure/cache"
	"github.com/ParkPawapon/request-api/internal/infrastructure/database"
	"github.com/ParkPawapon/request-api/internal/server"
	httptransport "github.com/ParkPawapon/request-api/internal/transport/http"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg         *config.Config
	db          *gorm.DB
	httpServer  *http.Server
	logger      *zap.Logger
	redisClient *redis.Client
}

func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*App, error) {
	db, err := database.Open(ctx, cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	redisClient := cache.NewClient(cfg.Redis)
	if err := cache.Ping(ctx, redisClient); err != nil {
		_ = database.Close(db)
		_ = cache.Close(redisClient)
		return nil, err
	}

	router := httptransport.NewRouter(httptransport.RouterDependencies{
		Config: cfg,
		DB:     db,
		Logger: logger,
		Redis:  redisClient,
	})

	return &App{
		cfg:         cfg,
		db:          db,
		httpServer:  server.NewHTTPServer(cfg.App, router),
		logger:      logger,
		redisClient: redisClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return server.RunWithGracefulShutdown(
		ctx,
		a.httpServer,
		func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
		},
		a.logger,
	)
}

func (a *App) Close(_ context.Context) {
	if err := cache.Close(a.redisClient); err != nil {
		a.logger.Warn("redis close failed", zap.Error(err))
	}
	if err := database.Close(a.db); err != nil {
		a.logger.Warn("database close failed", zap.Error(err))
	}
}
