package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ParkPawapon/request-api/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func Ping(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return client.Ping(pingCtx).Err()
}

func Close(client *redis.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}
