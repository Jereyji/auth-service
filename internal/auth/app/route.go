package auth_app

import (
	"github.com/Jereyji/auth-service/internal/auth/presentation/handlers"
	"github.com/gin-gonic/gin"
)


func InitRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshTokens)
		auth.POST("/dummyLogin", authHandler.DummyLogin)
	}
}
