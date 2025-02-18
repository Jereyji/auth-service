package handlers

import (
	"net/http"

	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

const (
	clientLevel    = 1
	moderatorLevel = 2
)

func (h Handler) AuthModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(h.accessTokenCookie.Name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, http.ErrNoCookie)
			return
		}

		claims, err := entity.ValidateAccessToken(accessToken, h.config.SecretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := claims.CheckAccessLevel(moderatorLevel); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Next()
	}
}
