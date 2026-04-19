package plantbatch

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"

	"github.com/gin-gonic/gin"
)

type activatePlantBatchItem struct {
	PackageID    uint    `json:"package_id" binding:"required,gt=0"`
	BatchCode    string  `json:"batch_code" binding:"required,min=3,max=60"`
	PlantingDate string  `json:"planting_date" binding:"required"`
	Location     string  `json:"location" binding:"required,min=3,max=150"`
	SeedCount    uint    `json:"seed_count" binding:"required,gt=0"`
	LandArea     float64 `json:"land_area" binding:"required,gt=0"`
}

type activatePackageRequest struct {
	PlantBatches []activatePlantBatchItem `json:"plant_batches" binding:"required"`
}

func ActivateUserPackage(c *gin.Context) {
	userID := c.Param("id")

	var req activatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant batch form"})
		return
	}
	if len(req.PlantBatches) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plant_batches is required"})
		return
	}

	var user authmodels.User
	if err := config.DB.Preload("Investments").Preload("SelectedPackage").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if len(user.Investments) == 0 && user.SelectedPackageID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has not selected a package yet"})
		return
	}

	now := time.Now()
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	onProcessInvestments := make([]coremodels.Investment, 0)
	for _, investment := range user.Investments {
		if investment.Status == "on_process" {
			onProcessInvestments = append(onProcessInvestments, investment)
		}
	}

	if len(user.Investments) == 0 && user.SelectedPackageID != nil {
		var packageItem coremodels.InvestmentPackage
		if err := tx.First(&packageItem, "id = ? AND status = ?", *user.SelectedPackageID, "active").Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Investment package not found"})
			return
		}

		legacyInvestment := coremodels.Investment{
			UserID:         userID,
			PackageID:      *user.SelectedPackageID,
			Amount:         packageItem.Price,
			Status:         "on_process",
			InvestmentDate: now,
		}

		if err := tx.Create(&legacyInvestment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
			return
		}

		onProcessInvestments = append(onProcessInvestments, legacyInvestment)
	}

	if len(onProcessInvestments) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "User must be in on_process status before activation"})
		return
	}

	if len(req.PlantBatches) != len(onProcessInvestments) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Each on_process package must have one plant batch form"})
		return
	}

	investmentIDs := make([]uint, 0, len(onProcessInvestments))
	onProcessByPackageID := make(map[uint]coremodels.Investment, len(onProcessInvestments))
	for _, investment := range onProcessInvestments {
		investmentIDs = append(investmentIDs, investment.ID)
		onProcessByPackageID[investment.PackageID] = investment
	}

	usedPackageIDs := map[uint]struct{}{}

	for _, item := range req.PlantBatches {
		_, exists := onProcessByPackageID[item.PackageID]
		if !exists {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Package %d is not in on_process status", item.PackageID)})
			return
		}

		if _, duplicate := usedPackageIDs[item.PackageID]; duplicate {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Duplicate plant batch for package %d", item.PackageID)})
			return
		}
		usedPackageIDs[item.PackageID] = struct{}{}
	}

	for _, item := range req.PlantBatches {
		investment := onProcessByPackageID[item.PackageID]

		plantingDate, err := time.Parse("2006-01-02", item.PlantingDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting_date format. Use YYYY-MM-DD"})
			return
		}

		batchCode := strings.TrimSpace(item.BatchCode)
		location := strings.TrimSpace(item.Location)
		if batchCode == "" || location == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "batch_code and location are required"})
			return
		}

		plantBatch := coremodels.PlantBatch{
			InvestmentID: investment.ID,
			BatchCode:    batchCode,
			PlantingDate: plantingDate,
			Location:     location,
			SeedCount:    item.SeedCount,
			LandArea:     item.LandArea,
			Status:       "active",
		}

		if err := tx.Create(&plantBatch).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save plant batch. Ensure batch_code is unique"})
			return
		}
	}

	if err := tx.Model(&coremodels.Investment{}).
		Where("id IN ?", investmentIDs).
		Updates(map[string]interface{}{
			"status":     "active",
			"updated_at": now,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	updates := map[string]interface{}{
		"package_status":       "active",
		"package_activated_at": now,
	}

	if err := tx.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package activated after plant batch validation"})
}