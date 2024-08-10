package api

import (
	_ "auth_service/api/docs"
	"auth_service/api/handler"
	"auth_service/pkg/logger"
	"auth_service/storage"

	"github.com/gin-gonic/gin"
	swaggerFile "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service
// @version 1.0
// @description This is the auth service of TravelTales app

// @contact.name Saidakbar
// @contact.url http://www.support_me_with_smile
// @contact.email "pardaboyevsaidakbar103@gmail.com"

// @host localhost:9999
// @BasePath /auth
func NewRouter(log logger.ILogger, storage storage.IStorage) *gin.Engine {
	r := gin.Default()

	r.GET("swagger/*any", swagger.WrapHandler(swaggerFile.Handler))

	h := handler.NewHandler(log, storage)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
	}

	return r
}
