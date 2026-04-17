package routes

import (
	"auth-app/handlers"
	"auth-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {

	api := r.Group("/api")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	user := api.Group("/users")
	user.Use(middleware.AuthMiddleware())
	user.GET("/", authHandler.GetUsers)
	user.GET("/:id", authHandler.GetUserByID)
	user.DELETE("/:id", authHandler.DeleteUser)
	user.PUT("/:id", authHandler.UpdateUser)
	
	utility := api.Group("/utility")
	utility.Use(middleware.AuthMiddleware())
	utility.GET("/quote", authHandler.GetTodayQuote)
}
