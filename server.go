package main

import (
	"synapsis/config"
	"synapsis/handler"
	"synapsis/middleware"
	"synapsis/repository"
	"synapsis/service"
	"synapsis/utils/http_response"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupServer(router *gin.Engine) {
	responseWriter := http_response.NewResponseWriter()
	database, err := config.NewDb(config.Envs.DB_URI)
	if err != nil {
		panic(err)
	}

	// repositories
	userRepo := repository.NewUserRepo(database)

	// services
	authService := service.NewAuthService(userRepo)

	// handlers
	authHandler := handler.NewAuthHandler(responseWriter, authService)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"error":   false,
			"message": "pong!",
		})
	})
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/check-token", authHandler.CheckToken)
	router.POST("/auth/register", authHandler.Register)

	secureRouter := router.Group("/")
	{
		secureRouter.Use(middleware.JWTMiddleware(responseWriter, authService))
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}