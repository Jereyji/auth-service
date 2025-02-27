package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jereyji/auth-service/internal/app/auth"
	"github.com/Jereyji/auth-service/internal/infrastucture/database/postgres"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
)

const (
	configPath = "config/config.yaml"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	config, err := configs.NewConfig(configPath)
	if err != nil {
		logger.Error("error reading environment variables: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	postgresDB, err := postgres.NewPostgresDB(ctx, config.Database)
	if err != nil {
		logger.Error("error initializing postgres database: ", slog.String("error", err.Error()))
	}
	defer postgresDB.Pool.Close()

	app := auth.NewAuthApp(ctx, config, logger, postgresDB)

	if err := app.Run(ctx); err != nil {
		logger.Error("server error: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
