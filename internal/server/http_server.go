package server

import (
	"fmt"
	"net/http"

	"github.com/ParkPawapon/request-api/internal/config"
)

func NewHTTPServer(cfg config.AppConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
