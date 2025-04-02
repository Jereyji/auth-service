package auth_app

import (
	"net/http"
	"strconv"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func (a AuthApp) AuthModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(accessTokenName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, http.ErrNoCookie)
			return
		}

		claims, err := entity.ValidateAccessToken(accessToken, a.secretMng.SecretKey)
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

func PrometheusMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		timer := prometheus.NewTimer(requestDuration.WithLabelValues(serviceName, path))

		c.Next()

		timer.ObserveDuration()

		status := c.Writer.Status()

		statusResponse.WithLabelValues(serviceName, path, strconv.Itoa(status)).Inc()
		totalRequests.WithLabelValues(serviceName, path).Inc()
	}
}
