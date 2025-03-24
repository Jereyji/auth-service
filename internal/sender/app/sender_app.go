package sender_app

import (
	"context"
	"log/slog"

	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/Jereyji/auth-service/internal/pkg/kafka"
	sender_service "github.com/Jereyji/auth-service/internal/sender/application/services"
	kafka_handler "github.com/Jereyji/auth-service/internal/sender/infrastructure/kafka"
)

type SenderApp struct {
	logger        *slog.Logger
	kafkaConsumer *kafka.KafkaConsumer
}

func NewSenderApp(ctx context.Context, config *configs.SenderConfig, logger *slog.Logger) (*SenderApp, error) {
	senderService := sender_service.NewSenderService(config)

	eventHandler := kafka_handler.NewLoginEventHandler(senderService, logger)

	consumer, err := kafka.NewKafkaConsumer(config.Kafka.Brokers, config.Kafka.GroupID, config.Kafka.Topic, eventHandler)
	if err != nil {
		return nil, err
	}

	return &SenderApp{
		kafkaConsumer: consumer,
		logger:        logger,
	}, nil
}

func (a *SenderApp) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		a.logger.Info("Received shutdown signal")

		if err := a.kafkaConsumer.Close(); err != nil {
			a.logger.Error("Error closing Kafka consumer", slog.String("error", err.Error()))
		}
	}()

	a.kafkaConsumer.StartListening(ctx)
	a.logger.Info("Sender service is running")

	<-ctx.Done()
	a.logger.Info("Sender service is shutting down")

	return nil
}
