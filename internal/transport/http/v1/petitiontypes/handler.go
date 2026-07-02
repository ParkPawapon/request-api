package petitiontypes

import (
	"context"
	"net/http"

	domain "github.com/ParkPawapon/request-api/internal/domain/petitiontype"
	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

type ListUseCase interface {
	Execute(ctx context.Context) ([]domain.PetitionType, error)
}

type Handler struct {
	listUseCase ListUseCase
}

func NewHandler(listUseCase ListUseCase) *Handler {
	return &Handler{listUseCase: listUseCase}
}

func (h *Handler) List(c *gin.Context) {
	if h == nil || h.listUseCase == nil {
		response.AppError(c, apperrors.ServiceNotReady("Petition type service is not ready.", nil))
		return
	}

	items, err := h.listUseCase.Execute(c.Request.Context())
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, toResponse(items))
}
