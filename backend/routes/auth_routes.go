package routes

import (
	authhandler "pondok-tani-backend/handlers/auth"
	userhandler "pondok-tani-backend/handlers/user"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authhandler.Register)
		auth.POST("/login", authhandler.Login)
	}
	auth.GET("/me", middleware.AuthMiddleware(), userhandler.GetCurrentUser)
}
