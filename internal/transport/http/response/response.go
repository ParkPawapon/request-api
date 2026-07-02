package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope[T any] struct {
	OK    bool       `json:"ok"`
	Data  *T         `json:"data,omitempty"`
	Meta  *Meta      `json:"meta,omitempty"`
	Error *ErrorBody `json:"error,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"pageSize,omitempty"`
	TotalItems int `json:"totalItems,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

func Success[T any](c *gin.Context, status int, data T) {
	c.JSON(status, Envelope[T]{
		OK:   true,
		Data: &data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
