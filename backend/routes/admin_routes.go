package routes

import (
	contenthandler "pondok-tani-backend/handlers/content"
	investmenthandler "pondok-tani-backend/handlers/investment"
	maintenancehandler "pondok-tani-backend/handlers/maintenance"
	notificationhandler "pondok-tani-backend/handlers/notification"
	plantbatchhandler "pondok-tani-backend/handlers/plantbatch"
	pricehandler "pondok-tani-backend/handlers/price"
	userhandler "pondok-tani-backend/handlers/user"
	dashboardhandler "pondok-tani-backend/handlers/dashboard"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	publicVanili := r.Group("/admin/vanili")
	{
		publicVanili.GET("/grades", pricehandler.ListGrades)
		publicVanili.GET("/prices", pricehandler.ListNationalPrices)
	}

	admin := r.Group("/admin", middleware.AuthMiddleware(), middleware.RequireRoles("admin"))
	{
		admin.GET("/dashboard", dashboardhandler.GetAdminDashboard)
		admin.GET("/users", userhandler.ListUsers)
		admin.GET("/users/:id", userhandler.GetUser)
		admin.POST("/users", userhandler.CreateUser)
		admin.PUT("/users/:id", userhandler.UpdateUser)
		admin.DELETE("/users/:id", userhandler.DeleteUser)
		admin.GET("/users/:id/plant-batches", plantbatchhandler.ListUserPlantBatches)
		admin.PUT("/users/:id/activate-package", plantbatchhandler.ActivateUserPackage)
		admin.PUT("/users/:id/investment-status", investmenthandler.UpdateInvestmentStatus)
		admin.POST("/notifications/push", notificationhandler.SendPushNotification)
		admin.GET("/emails/templates", notificationhandler.ListEmailTemplates)
		admin.POST("/emails/templates", notificationhandler.CreateEmailTemplate)
		admin.PUT("/emails/templates/:id", notificationhandler.UpdateEmailTemplate)
		admin.DELETE("/emails/templates/:id", notificationhandler.DeleteEmailTemplate)
		admin.GET("/emails/logs", notificationhandler.ListEmailLogs)
		admin.GET("/emails/monitoring", notificationhandler.ListEmailMonitoring)
		admin.POST("/emails/check-duplicate", notificationhandler.CheckEmailDuplicate)
		admin.POST("/emails/send", notificationhandler.SendEmailLog)
		admin.GET("/maintenance-schedules/summary", maintenancehandler.GetMaintenanceScheduleSummary)
		admin.GET("/maintenance-schedules", maintenancehandler.ListMaintenanceSchedules)
		admin.POST("/maintenance-schedules", maintenancehandler.CreateMaintenanceSchedule)
		admin.GET("/maintenance-activities/summary", maintenancehandler.GetAdminMaintenanceActivitySummary)
		admin.GET("/maintenance-activities", maintenancehandler.ListAdminMaintenanceActivities)
		admin.PUT("/maintenance-activities/:id/review", maintenancehandler.ReviewMaintenanceActivity)
		admin.GET("/articles", contenthandler.ListArticles)
		admin.GET("/articles/:id", contenthandler.GetArticle)
		admin.POST("/articles", contenthandler.CreateArticle)
		admin.PUT("/articles/:id", contenthandler.UpdateArticle)
		admin.DELETE("/articles/:id", contenthandler.DeleteArticle)

		admin.POST("/vanili/grades", pricehandler.CreateGrade)
		admin.PUT("/vanili/grades/:id", pricehandler.UpdateGrade)
		admin.DELETE("/vanili/grades/:id", pricehandler.DeleteGrade)
		admin.POST("/vanili/prices", pricehandler.CreateNationalPrice)
		admin.PUT("/vanili/prices/:id", pricehandler.UpdateNationalPrice)
		admin.DELETE("/vanili/prices/:id", pricehandler.DeleteNationalPrice)
	}
}
