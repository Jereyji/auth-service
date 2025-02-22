package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	repos "github.com/Jereyji/auth-service.git/internal/domain/interface_repository"
	"github.com/gin-gonic/gin"
)

func (h Handler) Register(c *gin.Context) {
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

		h.slog.Error("registration user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Email: user.Email,
	})
}

func (h Handler) Login(c *gin.Context) {
	var user LoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			c.String(http.StatusNotFound, "%s : %s", err.Error(), user.Email)
			return
		}

		if errors.Is(err, entity.ErrInvalidUsernameOrPassword) {
			c.String(http.StatusUnauthorized, entity.ErrInvalidUsernameOrPassword.Error())
			return
		}

		h.slog.Error("login user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)
	c.Status(http.StatusOK)
}

func (h Handler) DummyLogin(c *gin.Context) {
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
			c.String(http.StatusNotFound, "%s : %s", err.Error(), user.Email)
			return
		}

		h.slog.Error("dummy login user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)

	c.Status(http.StatusOK)
}

func (h Handler) RefreshTokens(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	accessToken, refreshToken, err := h.service.RefreshTokens(c.Request.Context(), refreshTokenCookie)
	if err != nil {
		// if errors.Is(err, repos.ErrNotFound) {
		// 	c.String(http.StatusNotFound, "%s : %s", err.Error(), refreshTokenCookie)
		// 	return
		// }

		h.slog.Error("refreshing tokens: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	h.sendTokensInCookie(c, accessToken.Token, refreshToken.RefreshToken)

	c.Status(http.StatusOK)
}

func (h Handler) Logout(c *gin.Context) {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.String(http.StatusBadRequest, ErrInvalidInput)
		return
	}

	if err := h.service.Logout(c.Request.Context(), refreshTokenCookie); err != nil {
		h.slog.Error("logouting user: ", slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
