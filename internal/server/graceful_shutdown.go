package server

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func RunWithGracefulShutdown(
	ctx context.Context,
	srv *http.Server,
	shutdownTimeoutContext func() (context.Context, context.CancelFunc),
	log *zap.Logger,
) error {
	errCh := make(chan error, 1)
	go func() {
		log.Info("http server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := shutdownTimeoutContext()
		defer cancel()
		log.Info("http server shutting down")
		return srv.Shutdown(shutdownCtx)
	}
}
