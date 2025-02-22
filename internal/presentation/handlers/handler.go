package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jereyji/auth-service.git/internal/application/services"
	"github.com/Jereyji/auth-service.git/internal/pkg/configs"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidInput = errors.New("invalid input").Error()
)

type Handler struct {
	service            *services.AuthService
	config             *configs.AuthConfig
	accessTokenCookie  http.Cookie
	refreshTokenCookie http.Cookie
	slog               *slog.Logger
}

func NewHandler(serv *services.AuthService, config *configs.AuthConfig, slog *slog.Logger) *Handler {
	return &Handler{
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
		slog: slog,
	}
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshTokens)
		auth.POST("/dummyLogin", h.DummyLogin)
	}

	return router
}

func (h Handler) sendTokensInCookie(c *gin.Context, accessToken, refreshToken string) {
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
