package health

import (
	"context"
	"fmt"
	"net/http"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckFunc func(ctx context.Context) error

type Dependencies struct {
	CheckDatabase CheckFunc
	CheckRedis    CheckFunc
	Logger        *zap.Logger
}

type Handler struct {
	checkDatabase CheckFunc
	checkRedis    CheckFunc
	logger        *zap.Logger
}

type Status struct {
	OK       bool   `json:"ok"`
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
	Redis    string `json:"redis,omitempty"`
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		checkDatabase: deps.CheckDatabase,
		checkRedis:    deps.CheckRedis,
		logger:        deps.Logger,
	}
}

func (h *Handler) Live(c *gin.Context) {
	response.Success(c, http.StatusOK, Status{
		OK:     true,
		Status: "live",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	dbStatus := "up"
	redisStatus := "up"
	var readinessErr error

	if h.checkDatabase == nil {
		dbStatus = "not_configured"
		readinessErr = fmt.Errorf("database checker is not configured")
	} else if err := h.checkDatabase(c.Request.Context()); err != nil {
		dbStatus = "down"
		readinessErr = err
	}

	if h.checkRedis == nil {
		redisStatus = "not_configured"
		if readinessErr == nil {
			readinessErr = fmt.Errorf("redis checker is not configured")
		}
	} else if err := h.checkRedis(c.Request.Context()); err != nil {
		redisStatus = "down"
		if readinessErr == nil {
			readinessErr = err
		}
	}

	if readinessErr != nil {
		if h.logger != nil {
			h.logger.Warn("readiness check failed", zap.Error(readinessErr))
		}
		response.AppError(c, apperrors.ServiceNotReady("Service is not ready.", readinessErr))
		return
	}

	response.Success(c, http.StatusOK, Status{
		OK:       true,
		Status:   "ready",
		Database: dbStatus,
		Redis:    redisStatus,
	})
}
