package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(appEnv string, level string) (*zap.Logger, error) {
	zapLevel := zapcore.InfoLevel
	if err := zapLevel.Set(strings.ToLower(strings.TrimSpace(level))); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	var cfg zap.Config
	if appEnv == "production" || appEnv == "staging" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	cfg.DisableStacktrace = appEnv == "production"

	return cfg.Build()
}

func Error(err error) zap.Field {
	return zap.Error(err)
}
