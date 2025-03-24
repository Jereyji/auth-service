package auth_app

import (
	"net/http"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	"github.com/gin-gonic/gin"
)

func (a AuthApp) AuthModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(a.accessTokenCookie.Name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, http.ErrNoCookie)
			return
		}

		claims, err := entity.ValidateAccessToken(accessToken, a.SecretMng.SecretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		userID := claims.TokenPayload.UserID
		c.Set("userID", userID)

		c.Next()
	}
}
