package config

import (
	"testing"
	"time"
)

func TestLoadValidConfig(t *testing.T) {
	t.Setenv("APP_NAME", "request-api")
	t.Setenv("APP_ENV", EnvTest)
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/request?sslmode=disable")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://request.example.test")
	t.Setenv("APP_REQUEST_TIMEOUT", "2s")
	t.Setenv("REDIS_KEY_PREFIX", "request-api:test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.App.Name != "request-api" {
		t.Fatalf("expected app name request-api, got %q", cfg.App.Name)
	}
	if cfg.App.RequestTimeout != 2*time.Second {
		t.Fatalf("expected request timeout 2s, got %s", cfg.App.RequestTimeout)
	}
	if len(cfg.CORS.AllowedOrigins) != 2 {
		t.Fatalf("expected 2 cors origins, got %d", len(cfg.CORS.AllowedOrigins))
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("APP_NAME", "request-api")
	t.Setenv("APP_ENV", EnvTest)
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("REDIS_ADDR", "localhost:6379")

	if _, err := Load(); err == nil {
		t.Fatal("expected missing DATABASE_URL to fail")
	}
}

func TestValidateRejectsInvalidCORSOrigin(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Name: "request-api",
			Env:  EnvTest,
			Port: "9090",
		},
		Database: DatabaseConfig{
			URL:          "postgres://user:pass@localhost:5432/request?sslmode=disable",
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		},
		Redis: RedisConfig{
			Addr: "localhost:6379",
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"localhost:3000"},
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid CORS origin to fail")
	}
}

func TestValidateRequiresCORSOriginsInProduction(t *testing.T) {
	cfg := validConfig()
	cfg.App.Env = EnvProduction
	cfg.CORS.AllowedOrigins = nil

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected production without CORS origins to fail")
	}
}

func TestValidateRejectsCORSOriginWithPath(t *testing.T) {
	cfg := validConfig()
	cfg.CORS.AllowedOrigins = []string{"https://request.example.test/api"}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected CORS origin with path to fail")
	}
}

func TestValidateRejectsInvalidPoolAndTimeouts(t *testing.T) {
	cfg := validConfig()
	cfg.App.RequestTimeout = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected zero request timeout to fail")
	}

	cfg = validConfig()
	cfg.Database.MaxIdleConns = cfg.Database.MaxOpenConns + 1

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected max idle conns greater than max open conns to fail")
	}
}

func TestValidateRejectsUnsafeRedisKeyPrefix(t *testing.T) {
	cfg := validConfig()
	cfg.Redis.KeyPrefix = "request api"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected unsafe redis key prefix to fail")
	}
}

func validConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:            "request-api",
			Env:             EnvTest,
			Port:            "9090",
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			ShutdownTimeout: time.Second,
			RequestTimeout:  time.Second,
			MaxBodyBytes:    1024,
		},
		Database: DatabaseConfig{
			URL:             "postgres://user:pass@localhost:5432/request?sslmode=disable",
			MaxOpenConns:    2,
			MaxIdleConns:    1,
			ConnMaxLifetime: time.Minute,
		},
		Redis: RedisConfig{
			Addr:      "localhost:6379",
			KeyPrefix: "request-api",
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"https://request.example.test"},
		},
	}
}
