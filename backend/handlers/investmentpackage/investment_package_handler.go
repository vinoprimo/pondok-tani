package investmentpackage

import (
	"net/http"
	"strings"

	"pondok-tani-backend/config"
	coremodels "pondok-tani-backend/models/core"

	"github.com/gin-gonic/gin"
)

type createOrUpdateInvestmentPackageRequest struct {
	PackageName string  `json:"package_name" binding:"required"`
	Description *string `json:"description"`
	MinQuantity uint    `json:"min_quantity" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Status      string  `json:"status"`
}

func normalizePackageStatus(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "active" || normalized == "inactive" {
		return normalized
	}
	return "active"
}

func ensureAdmin(c *gin.Context) bool {
	roleValue, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return false
	}

	role, ok := roleValue.(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can perform this action"})
		return false
	}

	return true
}

func ListInvestmentPackages(c *gin.Context) {
	var items []coremodels.InvestmentPackage

	query := config.DB.Order("id DESC")
	if status := strings.ToLower(strings.TrimSpace(c.Query("status"))); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch investment packages"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func GetInvestmentPackage(c *gin.Context) {
	var item coremodels.InvestmentPackage
	if err := config.DB.First(&item, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment package not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func CreateInvestmentPackage(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req createOrUpdateInvestmentPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}
	if req.MinQuantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_quantity must be greater than 0"})
		return
	}

	item := coremodels.InvestmentPackage{
		PackageName: strings.TrimSpace(req.PackageName),
		Description: req.Description,
		MinQuantity: req.MinQuantity,
		Price:       req.Price,
		Status:      normalizePackageStatus(req.Status),
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create investment package"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func UpdateInvestmentPackage(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var item coremodels.InvestmentPackage
	if err := config.DB.First(&item, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment package not found"})
		return
	}

	var req createOrUpdateInvestmentPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}
	if req.MinQuantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_quantity must be greater than 0"})
		return
	}

	item.PackageName = strings.TrimSpace(req.PackageName)
	item.Description = req.Description
	item.MinQuantity = req.MinQuantity
	item.Price = req.Price
	item.Status = normalizePackageStatus(req.Status)

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment package"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func DeleteInvestmentPackage(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	if err := config.DB.Delete(&coremodels.InvestmentPackage{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete investment package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment package deleted"})
}
