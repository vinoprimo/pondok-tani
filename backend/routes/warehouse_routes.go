package routes

import (
	warehousehandler "pondok-tani-backend/handlers/warehouse"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func WarehouseRoutes(r *gin.Engine) {
	group := r.Group("/warehouse", middleware.AuthMiddleware())
	{
		group.GET("/summary", warehousehandler.GetWarehouseSummary)
		group.GET("/stocks", warehousehandler.ListWarehouseStocks)
		group.GET("/movements", warehousehandler.ListStockMovements)
	}
}
