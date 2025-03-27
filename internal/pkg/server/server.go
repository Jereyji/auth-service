package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type HTTPServer struct {
	srv    *http.Server
	logger *slog.Logger
}

func NewHTTPServer(
	ctx context.Context,
	address string,
	handler http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	logger *slog.Logger,
) *HTTPServer {
	s := HTTPServer{
		srv: &http.Server{
			Addr:         address,
			Handler:      handler,
			BaseContext:  func(net.Listener) context.Context { return ctx },
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		logger: logger,
	}

	return &s
}

func (s *HTTPServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()

		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			s.logger.Warn("failed shutdown http server", slog.String("error", err.Error()))
		}
	}()

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
