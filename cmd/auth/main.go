package main

import (
	"context"
	"fmt"
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

	var config configs.AuthConfig
	if err := configs.NewConfig(&config, configPath, envPath); err != nil {
		panic(fmt.Errorf("error reading config variables: %w", err))
	}

	postgresDB, err := postgres.NewPostgresDB(ctx, &config.Postgres)
	if err != nil {
		panic(fmt.Errorf("error initializing postgres database: %w", err))
	}
	defer postgresDB.Pool.Close()

	kafkaProducer, err := kafka.NewKafkaProducer(config.Kafka.Brokers, config.Kafka.Topic)
	if err != nil {
		panic(fmt.Errorf("error initializing kafka producer: %w", err))
	}

	logger := SetupLogger(config.Gin.Mode)

	app := auth_app.NewAuthApp(ctx, &config, kafkaProducer, postgresDB, logger)

	if err := app.Run(ctx); err != nil {
		logger.Error("server error: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func SetupLogger(ginMode string) *slog.Logger {
	var logger *slog.Logger

	switch ginMode {
	case "debug":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "release":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
