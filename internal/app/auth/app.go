package auth

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Jereyji/auth-service/internal/application/services"
	"github.com/Jereyji/auth-service/internal/infrastucture/database/postgres"
	repository "github.com/Jereyji/auth-service/internal/infrastucture/repository/postgres"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/Jereyji/auth-service/internal/presentation/handlers"
	trm "github.com/Jereyji/auth-service/pkg/transaction_manager"
	"github.com/gin-gonic/gin"
)

type AuthApp struct {
	logger     *slog.Logger
	httpServer *http.Server
}

func NewAuthApp(
	ctx context.Context,
	cfg *configs.Config,
	logger *slog.Logger,
	postgresDB *postgres.PostgresDB,
) *AuthApp {
	trm := trm.NewTransactionManager(postgresDB.Pool)

	var (
		repos   = repository.NewAuthRepository(trm)
		service = services.NewAuthService(repos, trm, &cfg.AuthService)
		handler = handlers.NewAuthHandler(service, &cfg.AuthService, logger)
	)

	router := gin.Default()

	InitRoutes(router, handler)

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		BaseContext:  func(net.Listener) context.Context { return ctx },
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &AuthApp{
		logger:     logger,
		httpServer: srv,
	}
}

func (a AuthApp) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			a.logger.Warn("failed shutdown http server", slog.String("error", err.Error()))
		}
	}()

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
