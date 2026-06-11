package routes

import (
	financehandler "pondok-tani-backend/handlers/finance"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func FinancialRoutes(r *gin.Engine) {
	group := r.Group("/financial", middleware.AuthMiddleware())
	{
		group.GET("/operational-costs", financehandler.ListOperationalCosts)
		group.POST("/operational-costs", financehandler.CreateOperationalCost)
		group.PUT("/operational-costs/:id", financehandler.UpdateOperationalCost)
		group.DELETE("/operational-costs/:id", financehandler.DeleteOperationalCost)

		group.GET("/revenue-simulations", financehandler.ListRevenueSimulations)
		group.POST("/revenue-simulations", financehandler.CreateRevenueSimulation)
		group.PUT("/revenue-simulations/:id", financehandler.UpdateRevenueSimulation)
		group.DELETE("/revenue-simulations/:id", financehandler.DeleteRevenueSimulation)
		group.GET("/revenues", financehandler.ListActualRevenues)
	}
}
