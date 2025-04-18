package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/Jereyji/auth-service/internal/pkg/kafka"

	// kafka_models "github.com/Jereyji/auth-service/internal/pkg/kafka/models"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Register(ctx context.Context, name string, email string, password string) error
	DummyLogin(ctx context.Context, name string, email string, password string) (entity.AccessToken, entity.RefreshToken, error)
	Login(ctx context.Context, email string, password string) (entity.AccessToken, entity.RefreshToken, error)
	RefreshTokens(ctx context.Context, refreshToken string) (entity.AccessToken, entity.RefreshToken, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthHandler struct {
	service       IAuthService
	kafkaProducer *kafka.KafkaProducer
	cookies       *Cookies
	logger        *slog.Logger
}

func NewAuthHandler(
	service IAuthService,
	kafkaProducer *kafka.KafkaProducer,
	tokensCfg *configs.TokensConfig,
	slog *slog.Logger,
) *AuthHandler {
	accessTokenCookie, refreshTokenCookie := initCookies(
		int(tokensCfg.AccessTokenExpiresIn.Seconds()),
		int(tokensCfg.RefreshTokenExpiresIn.Seconds()),
	)

	return &AuthHandler{
		service:       service,
		kafkaProducer: kafkaProducer,
		cookies: &Cookies{
			AccessToken:  &accessTokenCookie,
			RefreshToken: &refreshTokenCookie,
		},
		logger: slog,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user RegisterRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, auth_errors.ErrInvalidJSONInput.Error())
		return
	}

	if err := h.service.Register(c.Request.Context(), user.Name, user.Email, user.Password); err != nil {
		if errors.Is(err, auth_errors.ErrRowExist) {
			c.String(http.StatusConflict, "%s: %s", err.Error(), user.Email)
			return
		}

		h.logger.Error("registration user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Email: user.Email,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var user LoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, auth_errors.ErrInvalidJSONInput.Error())
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, auth_errors.ErrNotFound) || errors.Is(err, auth_errors.ErrInvalidEmailOrPassword) {
			c.String(http.StatusUnauthorized, "%s : %s", auth_errors.ErrInvalidEmailOrPassword, user.Email)
			return
		}

		h.logger.Error("login user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	// if err := h.sendEventToKafka(user.Email, kafka_models.LoginEvent{}); err != nil {
	// 	h.logger.Error("sending event error: ", slog.String("error", err.Error()))
	// 	c.Status(http.StatusInternalServerError)
	// 	return
	// }

	h.cookies.sendTokens(c, accessToken.Token, refreshToken.Token)
	c.Status(http.StatusOK)
}

func (h *AuthHandler) DummyLogin(c *gin.Context) {
	var user RegisterRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.String(http.StatusBadRequest, auth_errors.ErrInvalidJSONInput.Error())
		return
	}

	accessToken, refreshToken, err := h.service.DummyLogin(c.Request.Context(), user.Email, user.Name, user.Password)
	if err != nil {
		if errors.Is(err, auth_errors.ErrRowExist) {
			c.String(http.StatusConflict, "%s: %s", err.Error(), user.Email)
			return
		}

		if errors.Is(err, auth_errors.ErrNotFound) {
			c.String(http.StatusUnauthorized, "%s : %s", err.Error(), user.Email)
			return
		}

		h.logger.Error("dummy login user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.cookies.sendTokens(c, accessToken.Token, refreshToken.Token)

	c.Status(http.StatusOK)
}

func (h *AuthHandler) RefreshTokens(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.String(http.StatusBadRequest, auth_errors.ErrInvalidJSONInput.Error())
		return
	}

	accessToken, refreshToken, err := h.service.RefreshTokens(c.Request.Context(), refreshTokenCookie)
	if err != nil {
		if errors.Is(err, auth_errors.ErrNotFound) {
			c.String(http.StatusUnauthorized, "%s : %s", err.Error(), refreshTokenCookie)
			return
		}

		h.logger.Error("refreshing tokens: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.cookies.sendTokens(c, accessToken.Token, refreshToken.Token)

	c.Status(http.StatusOK)
}

func (h *AuthHandler) sendEventToKafka(email string, event any) error {
	// loginEvent := models.LoginEvent{
	// 	Email:     email,
	// 	Timestamp: time.Now(),
	// 	Success:   true,
	// }

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err := h.kafkaProducer.SendMessage(email, string(eventJSON)); err != nil {
		return err
	}

	return nil
}
