package database

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

func NewGormLogger(log *zap.Logger) logger.Interface {
	if log == nil {
		return logger.Default.LogMode(logger.Warn)
	}
	return logger.New(
		zap.NewStdLog(log.Named("gorm")),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
}
