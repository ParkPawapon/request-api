package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := newInMemoryRateLimiter(maxRequests, window)
	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.Header("Retry-After", strconv.Itoa(int(limiter.retryAfter().Seconds())))
			response.AppError(c, apperrors.New(
				apperrors.CodeRateLimited,
				http.StatusTooManyRequests,
				"Too many requests.",
				nil,
			))
			return
		}
		c.Next()
	}
}

type inMemoryRateLimiter struct {
	maxRequests int
	window      time.Duration
	mu          sync.Mutex
	clients     map[string]rateLimitWindow
	now         func() time.Time
}

type rateLimitWindow struct {
	count   int
	expires time.Time
}

func newInMemoryRateLimiter(maxRequests int, window time.Duration) *inMemoryRateLimiter {
	if maxRequests < 1 {
		maxRequests = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &inMemoryRateLimiter{
		maxRequests: maxRequests,
		window:      window,
		clients:     make(map[string]rateLimitWindow),
		now:         time.Now,
	}
}

func (l *inMemoryRateLimiter) allow(key string) bool {
	if key == "" {
		key = "unknown"
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	for client, window := range l.clients {
		if now.After(window.expires) {
			delete(l.clients, client)
		}
	}

	window := l.clients[key]
	if window.expires.IsZero() || now.After(window.expires) {
		l.clients[key] = rateLimitWindow{
			count:   1,
			expires: now.Add(l.window),
		}
		return true
	}

	if window.count >= l.maxRequests {
		return false
	}

	window.count++
	l.clients[key] = window
	return true
}

func (l *inMemoryRateLimiter) retryAfter() time.Duration {
	return l.window
}
