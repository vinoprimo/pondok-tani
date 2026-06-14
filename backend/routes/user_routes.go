package routes

import (
	contenthandler "pondok-tani-backend/handlers/content"
	investmenthandler "pondok-tani-backend/handlers/investment"
	maintenancehandler "pondok-tani-backend/handlers/maintenance"
	notificationhandler "pondok-tani-backend/handlers/notification"
	userhandler "pondok-tani-backend/handlers/user"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/articles", contenthandler.ListArticles)
	r.GET("/articles/:id", contenthandler.GetArticle)

	user := r.Group("/users", middleware.AuthMiddleware())
	{
		user.GET("/me", userhandler.GetCurrentUser)
		user.PUT("/me/fcm-token", userhandler.UpdateFCMToken)
		user.PUT("/me/package", middleware.RequireRoles("investor", "mitra"), investmenthandler.SelectPackage)
		user.GET("/me/maintenance-schedules", middleware.RequireRoles("investor", "mitra"), maintenancehandler.ListMyMaintenanceSchedules)
		user.GET("/me/maintenance-activities", middleware.RequireRoles("investor", "mitra"), maintenancehandler.ListMyMaintenanceActivities)
		user.POST("/me/maintenance-activities", middleware.RequireRoles("investor", "mitra"), maintenancehandler.SubmitMyMaintenanceActivity)
		user.GET("/me/notifications", notificationhandler.GetMyNotifications)
		user.PUT("/me/notifications/:id/read", notificationhandler.MarkNotificationRead)
	}
}
