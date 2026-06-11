package routes

import (
	harvesthandler "pondok-tani-backend/handlers/harvest"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func HarvestRequestRoutes(r *gin.Engine) {
	handler := harvesthandler.NewHarvestRequestHandler()

	group := r.Group("/harvest-requests", middleware.AuthMiddleware())
	{
		group.GET("", handler.ListHarvestRequests)
		group.POST("", handler.CreateHarvestRequest)
		group.PATCH("/:id/validate", middleware.RequireRoles("admin"), handler.ValidateHarvestRequest)
	}
}
