package routes

import (
	userhandler "pondok-tani-backend/handlers/user"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/users", middleware.AuthMiddleware())
	{
		user.GET("/me", userhandler.GetCurrentUser)
		user.PUT("/me/package", middleware.RequireRoles("investor", "mitra"), userhandler.SelectPackage)
	}
}
