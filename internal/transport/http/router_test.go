package httptransport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ParkPawapon/request-api/internal/config"
	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"go.uber.org/zap"
)

func TestRouterNoRouteReturnsNormalizedError(t *testing.T) {
	router := NewRouter(RouterDependencies{
		Config: testConfig(),
		Logger: zap.NewNop(),
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/legacy-route", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var body response.Envelope[struct{}]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.OK || body.Error == nil {
		t.Fatalf("expected error envelope, got %+v", body)
	}
	if body.Error.Code != apperrors.CodeNotFound || body.Error.Message != "Not Found" {
		t.Fatalf("unexpected error body: %+v", body.Error)
	}
}

func testConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Name:            "request-api",
			Env:             config.EnvTest,
			Port:            "9090",
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			ShutdownTimeout: time.Second,
			RequestTimeout:  time.Second,
			MaxBodyBytes:    1024,
		},
		Database: config.DatabaseConfig{
			URL:             "postgres://user:pass@localhost:5432/request?sslmode=disable",
			MaxOpenConns:    2,
			MaxIdleConns:    1,
			ConnMaxLifetime: time.Minute,
		},
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			KeyPrefix: "request-api",
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"https://request.example.test"},
		},
	}
}
