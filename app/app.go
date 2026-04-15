package app

import (
	"auth-app/config"
	"auth-app/handlers"
	"auth-app/repositories"
	"auth-app/routes"
	"auth-app/services"

	"github.com/gin-gonic/gin"
)

func InitializeApp() *gin.Engine {
	config.ConnectDB()

	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)

	quoteRepo := repositories.NewQuoteRepository(config.DB)
	quoteService := services.NewQuoteService(quoteRepo)

	authHandler := handlers.NewAuthHandler(userService, quoteService)

	r := gin.Default()
	routes.SetupRoutes(r, authHandler)

	return r
}
