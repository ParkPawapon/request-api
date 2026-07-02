package response

import (
	"errors"
	"net/http"

	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/pkg/validation"
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code      apperrors.Code      `json:"code"`
	Message   string              `json:"message"`
	RequestID string              `json:"requestId,omitempty"`
	Fields    []ValidationFailure `json:"fields,omitempty"`
}

type ValidationFailure struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func AppError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		writeError(c, appErr.Status, ErrorBody{
			Code:      appErr.Code,
			Message:   appErr.Message,
			RequestID: requestID(c),
		})
		return
	}

	appErr = apperrors.Internal("Internal Server Error", err)
	writeError(c, http.StatusInternalServerError, ErrorBody{
		Code:      appErr.Code,
		Message:   appErr.Message,
		RequestID: requestID(c),
	})
}

func ValidationError(c *gin.Context, failures []ValidationFailure) {
	writeError(c, http.StatusBadRequest, ErrorBody{
		Code:      apperrors.CodeValidation,
		Message:   "Invalid request.",
		RequestID: requestID(c),
		Fields:    failures,
	})
}

func ValidationFailures(failures []validation.Failure) []ValidationFailure {
	out := make([]ValidationFailure, 0, len(failures))
	for _, failure := range failures {
		out = append(out, ValidationFailure{
			Field:   failure.Field,
			Message: failure.Message,
		})
	}
	return out
}

func writeError(c *gin.Context, status int, body ErrorBody) {
	c.AbortWithStatusJSON(status, Envelope[struct{}]{
		OK:    false,
		Error: &body,
	})
}

func requestID(c *gin.Context) string {
	value, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	id, ok := value.(string)
	if !ok {
		return ""
	}
	return id
}
