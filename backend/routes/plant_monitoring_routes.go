package routes

import (
	plantmonitoringhandler "pondok-tani-backend/handlers/plantmonitoring"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func PlantMonitoringRoutes(r *gin.Engine) {
	group := r.Group("/plant-monitorings", middleware.AuthMiddleware())
	{
		group.GET("", plantmonitoringhandler.ListPlantMonitorings)
		group.POST("", plantmonitoringhandler.CreatePlantMonitoring)
	}
}
