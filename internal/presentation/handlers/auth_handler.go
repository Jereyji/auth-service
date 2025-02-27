package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jereyji/auth-service/internal/application/services"
	"github.com/Jereyji/auth-service/internal/domain/entity"
	repos "github.com/Jereyji/auth-service/internal/domain/interface_repository"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidInput = errors.New("invalid input").Error()
)

type AuthHandler struct {
	service            *services.AuthService
	config             *configs.AuthConfig
	accessTokenCookie  http.Cookie
	refreshTokenCookie http.Cookie
	logger             *slog.Logger
}

func NewAuthHandler(serv *services.AuthService, config *configs.AuthConfig, slog *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: serv,
		config:  config,
		accessTokenCookie: http.Cookie{
			Name:     "access_token",
			Path:     "/",
			Domain:   "",
			MaxAge:   int(config.AccessTokenExpiresIn.Seconds()),
			Secure:   false,
			HttpOnly: true,
		},
		refreshTokenCookie: http.Cookie{
			Name:     "refresh_token",
			Path:     "/auth",
			Domain:   "",
			MaxAge:   int(config.RefreshTokenExpiresIn.Seconds()),
			Secure:   false,
			HttpOnly: true,
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

	h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)
	c.Status(http.StatusOK)
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
		h.accessTokenCookie.Name,
		accessToken,
		h.accessTokenCookie.MaxAge,
		h.accessTokenCookie.Path,
		h.accessTokenCookie.Domain,
		h.accessTokenCookie.Secure,
		h.accessTokenCookie.HttpOnly,
	)

	c.SetCookie(
		h.refreshTokenCookie.Name,
		refreshToken,
		h.refreshTokenCookie.MaxAge,
		h.refreshTokenCookie.Path,
		h.refreshTokenCookie.Domain,
		h.refreshTokenCookie.Secure,
		h.refreshTokenCookie.HttpOnly,
	)
}
