package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	repos "github.com/Jereyji/auth-service/internal/auth/domain/interface_repository"
	"github.com/Jereyji/auth-service/internal/pkg/kafka"
	"github.com/Jereyji/auth-service/internal/pkg/kafka/models"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidInput = errors.New("invalid input").Error()
)

type AuthServiceI interface {
	Register(ctx context.Context, name string, email string, password string) error
	DummyLogin(ctx context.Context, name string, email string, password string) (*entity.AccessToken, *entity.RefreshSessions, error)
	Login(ctx context.Context, email string, password string) (*entity.AccessToken, *entity.RefreshSessions, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*entity.AccessToken, *entity.RefreshSessions, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Cookies struct {
	AccessToken  *http.Cookie
	RefreshToken *http.Cookie
}

type AuthHandler struct {
	service       AuthServiceI
	kafkaProducer *kafka.KafkaProducer
	cookies       *Cookies
	logger        *slog.Logger
}

func NewAuthHandler(
	service AuthServiceI,
	kafkaProducer *kafka.KafkaProducer,
	accessTokenCookie *http.Cookie,
	refreshTokenCookie *http.Cookie,
	slog *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		service:       service,
		kafkaProducer: kafkaProducer,
		cookies: &Cookies{
			AccessToken:  accessTokenCookie,
			RefreshToken: refreshTokenCookie,
		},
		logger: slog,
	}
}

func (h AuthHandler) Register(c *gin.Context) {
	var user RegisterRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := h.service.Register(c.Request.Context(), user.Name, user.Email, user.Password); err != nil {
		if errors.Is(err, repos.ErrRowExist) {
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

func (h AuthHandler) Login(c *gin.Context) {
	var user LoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) || errors.Is(err, entity.ErrInvalidEmailOrPassword) {
			c.String(http.StatusUnauthorized, "%s : %s", entity.ErrInvalidEmailOrPassword, user.Email)
			return
		}

		h.logger.Error("login user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := h.sendEventToKafka(user.Email); err != nil {
		h.logger.Error("sending event error: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	// h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)
	// c.Status(http.StatusOK)
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.RefreshToken,
	})
}

func (h AuthHandler) DummyLogin(c *gin.Context) {
	var user RegisterRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	accessToken, refreshToken, err := h.service.DummyLogin(c.Request.Context(), user.Email, user.Name, user.Password)
	if err != nil {
		if errors.Is(err, repos.ErrRowExist) {
			c.String(http.StatusConflict, "%s: %s", err.Error(), user.Email)
			return
		}

		if errors.Is(err, repos.ErrNotFound) {
			c.String(http.StatusUnauthorized, "%s : %s", err.Error(), user.Email)
			return
		}

		h.logger.Error("dummy login user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)

	c.Status(http.StatusOK)
}

func (h AuthHandler) RefreshTokens(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	accessToken, refreshToken, err := h.service.RefreshTokens(c.Request.Context(), refreshTokenCookie)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			c.String(http.StatusUnauthorized, "%s : %s", err.Error(), refreshTokenCookie)
			return
		}

		h.logger.Error("refreshing tokens: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)

	c.Status(http.StatusOK)
}

func (h AuthHandler) sendTokensInCookie(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(
		h.cookies.AccessToken.Name,
		accessToken,
		h.cookies.AccessToken.MaxAge,
		h.cookies.AccessToken.Path,
		h.cookies.AccessToken.Domain,
		h.cookies.AccessToken.Secure,
		h.cookies.AccessToken.HttpOnly,
	)

	c.SetCookie(
		h.cookies.RefreshToken.Name,
		refreshToken,
		h.cookies.RefreshToken.MaxAge,
		h.cookies.RefreshToken.Path,
		h.cookies.RefreshToken.Domain,
		h.cookies.RefreshToken.Secure,
		h.cookies.RefreshToken.HttpOnly,
	)
}

func (h AuthHandler) sendEventToKafka(email string) error {
	loginEvent := models.LoginEvent{
		Email:     email,
		Timestamp: time.Now(),
		Success:   true,
	}

	eventJSON, err := json.Marshal(loginEvent)
	if err != nil {
		return err
	}

	if err := h.kafkaProducer.SendMessage(email, string(eventJSON)); err != nil {
		return err
	}

	return nil
}
