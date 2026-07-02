package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

func TestRateLimitRejectsRequestsOverLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/limited", RateLimit(2, time.Minute), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/limited", nil)
		req.RemoteAddr = "192.0.2.10:1234"
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("request %d: expected status %d, got %d", i+1, http.StatusNoContent, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req.RemoteAddr = "192.0.2.10:1234"
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
	if w.Header().Get("Retry-After") != "60" {
		t.Fatalf("expected Retry-After 60, got %q", w.Header().Get("Retry-After"))
	}

	var body response.Envelope[struct{}]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error == nil || body.Error.Code != apperrors.CodeRateLimited {
		t.Fatalf("expected rate limited error, got %+v", body)
	}
}
