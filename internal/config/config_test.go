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
