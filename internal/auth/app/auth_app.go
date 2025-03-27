package auth_app

import (
	"context"
	"log/slog"
	"net/http"

	auth_service "github.com/Jereyji/auth-service/internal/auth/application/services"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/postgres"
	repository "github.com/Jereyji/auth-service/internal/auth/infrastucture/repository/postgres"
	"github.com/Jereyji/auth-service/internal/auth/presentation/handlers"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/Jereyji/auth-service/internal/pkg/kafka"
	"github.com/Jereyji/auth-service/internal/pkg/server"
	trm "github.com/Jereyji/auth-service/internal/pkg/transaction_manager"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ServiceName = "Auth Service"
)

type SecretManager struct {
	SecretKey string
}

type AuthApp struct {
	httpServer        *server.HTTPServer
	accessTokenCookie *http.Cookie
	SecretMng         *SecretManager
	logger            *slog.Logger
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

	router.Use(
		PrometheusMiddleware(ServiceName),
	)

	prometheus.MustRegister(totalRequests, statusResponse, requestDuration)

	InitRoutes(router, handler)
	InitPrometheusRoutes(router)

	srv := server.NewHTTPServer(ctx, cfg.Server.Address, router, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, logger)

	return  &AuthApp{
		httpServer:        srv,
		accessTokenCookie: &accessTokenCookie,
		SecretMng:         &SecretManager{cfg.Application.Tokens.SecretKey},
		logger:            logger,
	}
}

func (a AuthApp) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
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
