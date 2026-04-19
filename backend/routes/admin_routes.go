package routes

import (
	investmenthandler "pondok-tani-backend/handlers/investment"
	plantbatchhandler "pondok-tani-backend/handlers/plantbatch"
	userhandler "pondok-tani-backend/handlers/user"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin", middleware.AuthMiddleware(), middleware.RequireRoles("admin"))
	{
		admin.GET("/users", userhandler.ListUsers)
		admin.POST("/users", userhandler.CreateUser)
		admin.PUT("/users/:id", userhandler.UpdateUser)
		admin.DELETE("/users/:id", userhandler.DeleteUser)
		admin.GET("/users/:id/plant-batches", plantbatchhandler.ListUserPlantBatches)
		admin.PUT("/users/:id/activate-package", plantbatchhandler.ActivateUserPackage)
		admin.PUT("/users/:id/investment-status", investmenthandler.UpdateInvestmentStatus)
	}
}
