package auth_app

import (
	"context"
	"log/slog"
	_ "net/http/pprof"

	// ginpprof "github.com/gin-contrib/pprof"

	auth_service "github.com/Jereyji/auth-service/internal/auth/application/services"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/postgres"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/redis"
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
	ServiceName     = "Auth Service"
	accessTokenName = "access_token"
)

type secretManager struct {
	SecretKey string
}

type AuthApp struct {
	httpServer *server.HTTPServer
	secretMng  *secretManager
	logger     *slog.Logger
}

func NewAuthApp(
	ctx context.Context,
	cfg *configs.AuthConfig,
	kafkaProducer *kafka.KafkaProducer,
	postgresDB *postgres.PostgresDB,
	logger *slog.Logger,
) *AuthApp {
	redisClient := redis.NewRedisClient(&cfg.Redis)
	trm := trm.NewTransactionManager(postgresDB.Pool)

	var (
		userRepos         = repository.NewAuthRepository(trm, redisClient)
		refreshTokenRepos = userRepos
		service           = auth_service.NewAuthService(trm, userRepos, refreshTokenRepos, &cfg.Tokens)
		handler           = handlers.NewAuthHandler(service, kafkaProducer, &cfg.Tokens, logger)
	)

	gin.SetMode(cfg.Gin.Mode)
	router := gin.New()

	// if gin.Mode() != gin.ReleaseMode {
	// 	ginpprof.Register(router)
	// }

	prometheus.MustRegister(
		totalRequests,
		statusResponse,
		requestDuration,
	)

	router.Use(
		gin.Recovery(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: cfg.Gin.SkipPaths,
		}),
		PrometheusMiddleware(ServiceName),
	)

	InitRoutes(router, handler)
	InitPrometheusRoutes(router)

	srv := server.NewHTTPServer(ctx, cfg.Server.Address, router, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, logger)
	app := AuthApp{
		httpServer: srv,
		secretMng:  &secretManager{cfg.Tokens.SecretKey},
		logger:     logger,
	}

	return &app
}

func (a *AuthApp) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}
