package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Cookies struct {
	AccessToken  *http.Cookie
	RefreshToken *http.Cookie
}

func initCookies(accessTokenExpiresIn, refreshTokenExpiresIn int) (http.Cookie, http.Cookie) {
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Domain:   "",
		MaxAge:   accessTokenExpiresIn,
		Secure:   true,
		HttpOnly: true,
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Path:     "/auth",
		Domain:   "",
		MaxAge:   refreshTokenExpiresIn,
		Secure:   true,
		HttpOnly: true,
	}

	return accessTokenCookie, refreshTokenCookie
}

func (k *Cookies) sendTokens(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(
		k.AccessToken.Name,
		accessToken,
		k.AccessToken.MaxAge,
		k.AccessToken.Path,
		k.AccessToken.Domain,
		k.AccessToken.Secure,
		k.AccessToken.HttpOnly,
	)

	c.SetCookie(
		k.RefreshToken.Name,
		refreshToken,
		k.RefreshToken.MaxAge,
		k.RefreshToken.Path,
		k.RefreshToken.Domain,
		k.RefreshToken.Secure,
		k.RefreshToken.HttpOnly,
	)
}
