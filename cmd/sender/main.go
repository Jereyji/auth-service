package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jereyji/auth-service/internal/pkg/configs"
	sender_app "github.com/Jereyji/auth-service/internal/sender/app"
)

const (
	configPath = string("config/sender_config.yaml")
	envPath    = string("deployments/.env")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var config configs.SenderConfig
	if err := configs.NewConfig(&config, configPath, envPath); err != nil {
		panic(fmt.Errorf("error reading environment variables: %w", err))
	}

	logger := SetupLogger(config.EnvMode)

	senderApp, err := sender_app.NewSenderApp(ctx, &config, logger)
	if err != nil {
		panic(fmt.Errorf("error initializing sender app: %w", err))
	}

	if err := senderApp.Run(ctx); err != nil {
		panic(fmt.Errorf("error while running sender app: %w", err))
	}
}

func SetupLogger(envMode string) *slog.Logger {
	var logger *slog.Logger

	switch envMode {
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
