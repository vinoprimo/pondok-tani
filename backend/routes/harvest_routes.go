package routes

import (
	"pondok-tani-backend/config"
	harvesthandler "pondok-tani-backend/handlers/harvest"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func HarvestRoutes(r *gin.Engine) {
	handler := harvesthandler.NewHarvestHandler(config.DB)

	harvest := r.Group("/harvest", middleware.AuthMiddleware())
	{
		harvest.GET("", handler.ListHarvests)
		harvest.GET("/summary", handler.GetHarvestSummary)
		harvest.POST("", handler.CreateHarvest)
		harvest.POST("/:id/outputs", handler.CreateHarvestOutput)
	}
}
