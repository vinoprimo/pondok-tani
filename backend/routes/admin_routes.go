package routes

import (
	investmenthandler "pondok-tani-backend/handlers/investment"
	maintenancehandler "pondok-tani-backend/handlers/maintenance"
	plantbatchhandler "pondok-tani-backend/handlers/plantbatch"
	pricehandler "pondok-tani-backend/handlers/price"
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
		admin.GET("/maintenance-schedules/summary", maintenancehandler.GetMaintenanceScheduleSummary)
		admin.GET("/maintenance-schedules", maintenancehandler.ListMaintenanceSchedules)
		admin.POST("/maintenance-schedules", maintenancehandler.CreateMaintenanceSchedule)

		admin.GET("/vanili/grades", pricehandler.ListGrades)
		admin.POST("/vanili/grades", pricehandler.CreateGrade)
		admin.PUT("/vanili/grades/:id", pricehandler.UpdateGrade)
		admin.DELETE("/vanili/grades/:id", pricehandler.DeleteGrade)
		admin.GET("/vanili/prices", pricehandler.ListNationalPrices)
		admin.POST("/vanili/prices", pricehandler.CreateNationalPrice)
		admin.PUT("/vanili/prices/:id", pricehandler.UpdateNationalPrice)
		admin.DELETE("/vanili/prices/:id", pricehandler.DeleteNationalPrice)
	}
}
