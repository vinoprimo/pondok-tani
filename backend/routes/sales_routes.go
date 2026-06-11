package routes

import (
	saleshandler "pondok-tani-backend/handlers/sales"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SalesRoutes(r *gin.Engine) {
	group := r.Group("/sales", middleware.AuthMiddleware())
	{
		group.GET("/orders", saleshandler.ListSalesOrders)
		group.POST("/orders", middleware.RequireRoles("admin"), saleshandler.CreateSalesOrder)
	}
}
