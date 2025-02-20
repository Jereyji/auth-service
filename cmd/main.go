package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jereyji/auth-service.git/internal/application/services"
	"github.com/Jereyji/auth-service.git/internal/infrastucture/database/postgres"
	repository "github.com/Jereyji/auth-service.git/internal/infrastucture/repository/postgres"
	"github.com/Jereyji/auth-service.git/internal/pkg/configs"
	"github.com/Jereyji/auth-service.git/internal/presentation/handlers"
	trm "github.com/Jereyji/auth-service.git/pkg/transaction_manager"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	config, err := configs.NewConfig()
	if err != nil {
		log.Fatal("Error reading environment variables: ", err)
	}

	postgresDB, err := postgres.NewPostgresDB(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatal("Error initialization postgres database: ", err)
	}
	defer postgresDB.Close()

	trm := trm.NewTransactionManager(postgresDB)

	var (
		repos   = repository.NewAuthRepository(trm)
		service = services.NewAuthService(repos, trm, &config.AuthService)
		handler = handlers.NewHandler(service, &config.AuthService, logger)
	)

	r := handler.InitRoutes()

	err = r.Run("0.0.0.0:8080")
	if err == nil {
		log.Fatal("Error running auth-service: ", err)
	}
}
