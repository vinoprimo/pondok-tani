package routes

import (
	"pondok-tani-backend/handlers"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func InvestmentRoutes(r *gin.Engine) {
	inv := r.Group("/investments")
	inv.Use(middleware.AuthMiddleware())
	{
		inv.POST("/", handlers.CreateInvestment)
	}
	inv.GET("/", handlers.GetMyInvestments)
	inv.GET("/:id", handlers.GetInvestmentDetail)
}