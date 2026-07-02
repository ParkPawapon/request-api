package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

func TestLiveEndpoint(t *testing.T) {
	router := testRouter(NewHandler(Dependencies{}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/health/live", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body response.Envelope[Status]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.OK || body.Data == nil || body.Data.Status != "live" {
		t.Fatalf("unexpected response body: %+v", body)
	}
}

func TestReadyEndpointSuccess(t *testing.T) {
	router := testRouter(NewHandler(Dependencies{
		CheckDatabase: func(context.Context) error { return nil },
		CheckRedis:    func(context.Context) error { return nil },
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/health/ready", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body response.Envelope[Status]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.OK || body.Data == nil || body.Data.Database != "up" || body.Data.Redis != "up" {
		t.Fatalf("unexpected response body: %+v", body)
	}
}

func TestReadyEndpointFailureIsNormalized(t *testing.T) {
	router := testRouter(NewHandler(Dependencies{
		CheckDatabase: func(context.Context) error {
			return errors.New("postgres password leaked detail")
		},
		CheckRedis: func(context.Context) error { return nil },
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/health/ready", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var body response.Envelope[Status]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.OK || body.Error == nil {
		t.Fatalf("expected error response, got %+v", body)
	}
	if body.Error.Code != apperrors.CodeServiceNotReady {
		t.Fatalf("expected code %s, got %s", apperrors.CodeServiceNotReady, body.Error.Code)
	}
	if body.Error.Message != "Service is not ready." {
		t.Fatalf("unexpected error message %q", body.Error.Message)
	}
}

func testRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/v1")
	RegisterRoutes(group, handler)
	return router
}
