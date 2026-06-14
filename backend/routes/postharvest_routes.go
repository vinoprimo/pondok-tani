package routes

import (
	postharvesthandler "pondok-tani-backend/handlers/postharvest"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func PostHarvestRoutes(r *gin.Engine) {
	drying := r.Group("/drying-processes", middleware.AuthMiddleware())
	{
		drying.GET("", postharvesthandler.ListDryingProcesses)
		drying.PATCH("/:id/complete", middleware.RequireRoles("admin"), postharvesthandler.CompleteDryingProcess)
	}

	grading := r.Group("/grading-batches", middleware.AuthMiddleware())
	{
		grading.GET("", postharvesthandler.ListGradingBatches)
		grading.POST("", middleware.RequireRoles("admin"), postharvesthandler.CreateGradingBatch)
	}
}
