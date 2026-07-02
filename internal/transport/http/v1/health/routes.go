package health

import "github.com/gin-gonic/gin"

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	router.GET("/health", handler.Ready)
	router.GET("/health/live", handler.Live)
	router.GET("/health/ready", handler.Ready)
}
