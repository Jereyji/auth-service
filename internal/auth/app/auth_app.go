package auth_app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	auth_service "github.com/Jereyji/auth-service/internal/auth/application/services"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/postgres"
	repository "github.com/Jereyji/auth-service/internal/auth/infrastucture/repository/postgres"
	"github.com/Jereyji/auth-service/internal/auth/presentation/handlers"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/Jereyji/auth-service/internal/pkg/kafka"
	trm "github.com/Jereyji/auth-service/internal/pkg/transaction_manager"
	"github.com/gin-gonic/gin"
)

type SecretManager struct {
	SecretKey string
}

type AuthApp struct {
	logger            *slog.Logger
	httpServer        *http.Server
	accessTokenCookie *http.Cookie
	SecretMng         *SecretManager
}

func NewAuthApp(
	ctx context.Context,
	cfg *configs.AuthConfig,
	kafkaProducer *kafka.KafkaProducer,
	postgresDB *postgres.PostgresDB,
	logger *slog.Logger,
) *AuthApp {
	trm := trm.NewTransactionManager(postgresDB.Pool)

	accessTokenCookie, refreshTokenCookie := initCookies(
		int(cfg.Application.Tokens.AccessTokenExpiresIn.Seconds()),
		int(cfg.Application.Tokens.RefreshTokenExpiresIn.Seconds()),
	)

	var (
		repos   = repository.NewAuthRepository(trm)
		service = auth_service.NewAuthService(repos, trm, &cfg.Application)
		handler = handlers.NewAuthHandler(service, kafkaProducer, &accessTokenCookie, &refreshTokenCookie, logger)
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
		httpServer:        srv,
		accessTokenCookie: &accessTokenCookie,
		SecretMng:         &SecretManager{cfg.Application.Tokens.SecretKey},
		logger:            logger,
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

func initCookies(accessTokenExpiresIn, refreshTokenExpiresIn int) (http.Cookie, http.Cookie) {
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Domain:   "",
		MaxAge:   accessTokenExpiresIn,
		Secure:   true,
		HttpOnly: true,
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Path:     "/auth",
		Domain:   "",
		MaxAge:   refreshTokenExpiresIn,
		Secure:   true,
		HttpOnly: true,
	}

	return accessTokenCookie, refreshTokenCookie
}
