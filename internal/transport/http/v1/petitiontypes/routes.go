package petitiontypes

import (
	"time"

	"github.com/ParkPawapon/request-api/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	router.GET("/petition-types", middleware.RateLimit(60, time.Minute), handler.List)
}
