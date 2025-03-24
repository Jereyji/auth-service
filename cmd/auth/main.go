package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	auth_app "github.com/Jereyji/auth-service/internal/auth/app"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/postgres"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/Jereyji/auth-service/internal/pkg/kafka"
)

const (
	configPath = string("config/auth_config.yaml")
	envPath    = string("deployments/.env")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var config configs.AuthConfig
	if err := configs.NewConfig(&config, configPath, envPath); err != nil {
		logger.Error("error reading environment variables: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	postgresDB, err := postgres.NewPostgresDB(ctx, config.Database)
	if err != nil {
		logger.Error("error initializing postgres database: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer postgresDB.Pool.Close()

	kafkaProducer, err := kafka.NewKafkaProducer(config.Kafka.Brokers, config.Kafka.Topic)
	if err != nil {
		logger.Error("error initializing kafka producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	app := auth_app.NewAuthApp(ctx, &config, kafkaProducer, postgresDB, logger)

	if err := app.Run(ctx); err != nil {
		logger.Error("server error: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
