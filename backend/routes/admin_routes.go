package routes

import (
	userhandler "pondok-tani-backend/handlers/user"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin", middleware.AuthMiddleware(), middleware.RequireRoles("admin"))
	{
		admin.GET("/users", userhandler.ListUsers)
		admin.PUT("/users/:id/activate-package", userhandler.ActivateUserPackage)
	}
}
