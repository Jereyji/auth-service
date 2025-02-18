package handlers

import (
	"errors"
	"net/http"

	"github.com/Jereyji/auth-service.git/internal/application/services"
	"github.com/Jereyji/auth-service.git/internal/pkg/configs"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service            *services.Service
	config             *configs.AuthConfig
	accessTokenCookie  http.Cookie
	refreshTokenCookie http.Cookie
}

func NewHandler(serv *services.Service, config *configs.AuthConfig) *Handler {
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

var (
	ErrInvalidInput = errors.New("invalid input").Error()
)
