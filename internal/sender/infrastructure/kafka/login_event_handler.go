package kafka_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/Jereyji/auth-service/internal/pkg/kafka/models"
	sender_service "github.com/Jereyji/auth-service/internal/sender/application/services"
)

type LoginEventHandler struct {
	messageService sender_service.MessageServiceI
	logger         *slog.Logger
}

func NewLoginEventHandler(messageService sender_service.MessageServiceI, logger *slog.Logger) *LoginEventHandler {
	return &LoginEventHandler{
		messageService: messageService,
		logger:         logger,
	}
}

func (h *LoginEventHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *LoginEventHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *LoginEventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.logger.Info("message received",
			slog.String("topic", message.Topic),
			slog.Int("partition", int(message.Partition)),
			slog.Int64("offset", message.Offset),
			slog.String("key", string(message.Key)))

		var loginEvent models.LoginEvent
		if err := json.Unmarshal(message.Value, &loginEvent); err != nil {
			h.logger.Error("message deserialization error", slog.String("error", err.Error()))
			session.MarkMessage(message, "")
			continue
		}

		if err := h.handleLoginEvent(context.Background(), loginEvent); err != nil {
			h.logger.Error("login event processing error", slog.String("error", err.Error()))
		} else {
			h.logger.Info("login event successfully processed",
				slog.String("email", loginEvent.Email),
				slog.Bool("success", loginEvent.Success))
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func (h *LoginEventHandler) handleLoginEvent(ctx context.Context, event models.LoginEvent) error {
	var subject, content string

	timeFormatted := event.Timestamp.Format("02.01.2006 15:04:05")

	if event.Success {
		subject = "Успешный вход в аккаунт"
		content = fmt.Sprintf(
			"Здравствуйте,\n\nМы обнаружили успешный вход в ваш аккаунт.\n\n"+
				"Время: %s\n"+
				"Если это были вы, то никаких действий не требуется.\n"+
				"Если это были не вы, пожалуйста, немедленно смените пароль и свяжитесь с поддержкой.", timeFormatted)
	} else {
		subject = "Неудачная попытка входа в аккаунт"
		content = fmt.Sprintf(
			"Здравствуйте,\n\nМы обнаружили неудачную попытку входа в ваш аккаунт.\n\n"+
				"Время: %s\n"+
				"Если это были вы, возможно, вы ввели неверный пароль. Попробуйте снова.\n"+
				"Если это были не вы, пожалуйста, убедитесь, что ваш пароль надежен и свяжитесь с поддержкой.", timeFormatted)
	}

	return h.messageService.SendEmail(ctx, event.Email, subject, content)
}
