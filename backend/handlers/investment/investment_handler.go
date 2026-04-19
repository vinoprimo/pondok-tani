package investment

import (
	"net/http"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

	updates := map[string]interface{}{
		"package_status": req.Status,
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