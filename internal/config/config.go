package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	CORS     CORSConfig
	Log      LogConfig
}

type AppConfig struct {
	Name            string
	Env             string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	MaxBodyBytes    int64
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type CORSConfig struct {
	AllowedOrigins []string
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:            getString("APP_NAME", "request-api"),
			Env:             getString("APP_ENV", EnvDevelopment),
			Port:            getString("APP_PORT", "8080"),
			ReadTimeout:     getDuration("APP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDuration("APP_WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout: getDuration("APP_SHUTDOWN_TIMEOUT", 10*time.Second),
			RequestTimeout:  getDuration("APP_REQUEST_TIMEOUT", 15*time.Second),
			MaxBodyBytes:    getInt64("APP_MAX_BODY_BYTES", 50<<20),
		},
		Database: DatabaseConfig{
			URL:             getString("DATABASE_URL", ""),
			MaxOpenConns:    getInt("DATABASE_MAX_OPEN_CONNS", 15),
			MaxIdleConns:    getInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDuration("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getString("REDIS_ADDR", ""),
			Password: getString("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins: getCSV("CORS_ALLOWED_ORIGINS"),
		},
		Log: LogConfig{
			Level: getString("LOG_LEVEL", "info"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is required")
	}
	if c.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if !isAllowedEnv(c.App.Env) {
		return fmt.Errorf("APP_ENV must be one of development, test, staging, production")
	}
	if _, err := strconv.Atoi(c.App.Port); err != nil {
		return fmt.Errorf("APP_PORT must be numeric")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Database.MaxOpenConns < 1 {
		return fmt.Errorf("DATABASE_MAX_OPEN_CONNS must be greater than zero")
	}
	if c.Database.MaxIdleConns < 0 {
		return fmt.Errorf("DATABASE_MAX_IDLE_CONNS must not be negative")
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}
	if _, _, err := net.SplitHostPort(c.Redis.Addr); err != nil {
		return fmt.Errorf("REDIS_ADDR must be host:port")
	}
	for _, origin := range c.CORS.AllowedOrigins {
		if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			return fmt.Errorf("CORS_ALLOWED_ORIGINS must contain absolute http(s) origins")
		}
	}
	return nil
}

func isAllowedEnv(env string) bool {
	switch env {
	case EnvDevelopment, EnvTest, EnvStaging, EnvProduction:
		return true
	default:
		return false
	}
}
