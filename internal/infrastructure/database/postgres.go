package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ParkPawapon/request-api/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(ctx context.Context, cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{
		Logger: NewGormLogger(log),
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("extract sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := Ping(ctx, db); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("postgres client is nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return sqlDB.PingContext(pingCtx)
}

func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func WithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if db == nil {
		return fmt.Errorf("postgres client is nil")
	}
	return db.WithContext(ctx).Transaction(fn, &sql.TxOptions{})
}
