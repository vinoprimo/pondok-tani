package routes

import (
	investmentpackagehandler "pondok-tani-backend/handlers/investmentpackage"
	"pondok-tani-backend/middleware"

	"github.com/gin-gonic/gin"
)

func InvestmentPackageRoutes(r *gin.Engine) {
	packages := r.Group("/investment-packages")
	{
		packages.GET("", investmentpackagehandler.ListInvestmentPackages)
		packages.GET("/:id", investmentpackagehandler.GetInvestmentPackage)
		packages.POST("", middleware.AuthMiddleware(), investmentpackagehandler.CreateInvestmentPackage)
		packages.PUT("/:id", middleware.AuthMiddleware(), investmentpackagehandler.UpdateInvestmentPackage)
		packages.DELETE("/:id", middleware.AuthMiddleware(), investmentpackagehandler.DeleteInvestmentPackage)
	}
}
