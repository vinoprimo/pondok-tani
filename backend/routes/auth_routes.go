package routes

import (
	"pondok-tani-backend/handlers"
	"pondok-tani-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}
	auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
}