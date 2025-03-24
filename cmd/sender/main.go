package main

import (
	"context"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var config configs.SenderConfig
	if err := configs.NewConfig(&config, configPath, envPath); err != nil {
		logger.Error("error reading environment variables: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	senderApp, err := sender_app.NewSenderApp(ctx, &config, logger)
	if err != nil {
		logger.Error("error initializing sender app:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := senderApp.Run(ctx); err != nil {
		logger.Error("error while running sender app", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
