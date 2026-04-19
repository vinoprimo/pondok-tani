package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user authmodels.User

	err := config.DB.
		Preload("SelectedPackage").
		Preload("Investments").
		Preload("Investments.Package").
		First(&user, "id = ?", userID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type packageSelectionRequest struct {
	PackageID  uint   `json:"package_id"`
	PackageIDs []uint `json:"package_ids"`
}

func SelectPackage(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req packageSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	packageIDs := make([]uint, 0, len(req.PackageIDs)+1)
	seen := map[uint]struct{}{}
	if req.PackageID > 0 {
		packageIDs = append(packageIDs, req.PackageID)
		seen[req.PackageID] = struct{}{}
	}
	for _, packageID := range req.PackageIDs {
		if packageID == 0 {
			continue
		}
		if _, exists := seen[packageID]; exists {
			continue
		}
		packageIDs = append(packageIDs, packageID)
		seen[packageID] = struct{}{}
	}

	if len(packageIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one package must be selected"})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	now := time.Now()
	primaryPackageID := packageIDs[0]

	for _, packageID := range packageIDs {
		var packageItem coremodels.InvestmentPackage
		if err := tx.First(&packageItem, "id = ? AND status = ?", packageID, "active").Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Investment package not found"})
			return
		}

		investment := coremodels.Investment{
			UserID:         userID,
			PackageID:      packageID,
			Amount:         packageItem.Price,
			Status:         "pending",
			InvestmentDate: now,
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "package_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"amount",
				"status",
				"investment_date",
				"updated_at",
			}),
		}).Create(&investment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save package selection"})
			return
		}
	}

	updates := map[string]interface{}{
		"selected_package_id":  primaryPackageID,
		"package_status":       "pending",
		"package_selected_at":  now,
		"package_activated_at": gorm.Expr("NULL"),
	}

	if err := tx.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save package selection"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save package selection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Package selection saved",
		"package_ids":    packageIDs,
		"package_id":     primaryPackageID,
		"packages_count": len(packageIDs),
	})
}

func ListUsers(c *gin.Context) {
	var users []authmodels.User
	if err := config.DB.
		Preload("SelectedPackage").
		Preload("Investments").
		Preload("Investments.Package").
		Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

type updateInvestmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=on_process"`
}

func UpdateInvestmentStatus(c *gin.Context) {
	userID := c.Param("id")

	var req updateInvestmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
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
			Status:         req.Status,
			InvestmentDate: now,
		}

		if err := tx.Create(&legacyInvestment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment status"})
			return
		}
	} else {
		if err := tx.Model(&coremodels.Investment{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"status":     req.Status,
				"updated_at": now,
			}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment status"})
			return
		}
	}

	var packageStatusValue string
	if req.Status == "on_process" {
		packageStatusValue = "on_process"
	} else if req.Status == "active" {
		packageStatusValue = "active"
	}

	updates := map[string]interface{}{
		"package_status": packageStatusValue,
	}

	if req.Status == "active" {
		updates["package_activated_at"] = now
	}

	if err := tx.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment status"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Investment status updated",
		"status":  req.Status,
	})
}

func ActivateUserPackage(c *gin.Context) {
	userID := c.Param("id")

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
