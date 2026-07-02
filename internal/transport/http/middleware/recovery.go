package middleware

import (
	"fmt"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if log != nil {
					log.Error("panic recovered",
						zap.Any("panic", recovered),
						zap.String("request_id", requestID(c)),
					)
				}
				response.AppError(c, apperrors.Internal("Internal Server Error", fmt.Errorf("panic recovered")))
			}
		}()
		c.Next()
	}
}
