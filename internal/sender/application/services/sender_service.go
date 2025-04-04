package sender_service

import (
	"context"
	"fmt"

	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"gopkg.in/gomail.v2"
)

type MessageServiceI interface {
	SendEmail(ctx context.Context, recipient, subject, content string) error
}

type SenderService struct {
	username string
	dialer   *gomail.Dialer
}

func NewSenderService(config *configs.SenderConfig) MessageServiceI {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.Username, config.Password)
	return &SenderService{
		username: config.Username,
		dialer:   dialer,
	}
}

func (s *SenderService) SendEmail(ctx context.Context, recipient, subject, content string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.username)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("fail send message to %s: %w", recipient, err)
	}

	return nil
}
