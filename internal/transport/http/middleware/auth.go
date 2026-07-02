package middleware

import (
	"context"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

type Principal struct {
	Subject string
	Roles   []string
}

type Authenticator interface {
	Authenticate(ctx context.Context, c *gin.Context) (Principal, error)
}

func RequireAuthenticated(auth Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if auth == nil {
			response.AppError(c, apperrors.ServiceNotReady("Authentication is not configured.", nil))
			return
		}

		principal, err := auth.Authenticate(c.Request.Context(), c)
		if err != nil {
			response.AppError(c, apperrors.Unauthorized("Authentication is required.", err))
			return
		}

		c.Set("principal", principal)
		c.Next()
	}
}

func CurrentPrincipal(c *gin.Context) (Principal, bool) {
	value, exists := c.Get("principal")
	if !exists {
		return Principal{}, false
	}
	principal, ok := value.(Principal)
	return principal, ok
}
