package user

import (
	"net/http"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user authmodels.User

	err := config.DB.Preload("SelectedPackage").First(&user, "id = ?", userID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type packageSelectionRequest struct {
	PackageID uint `json:"package_id" binding:"required"`
}

func SelectPackage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req packageSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var packageItem coremodels.InvestmentPackage
	if err := config.DB.First(&packageItem, "id = ? AND status = ?", req.PackageID, "active").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment package not found"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"selected_package_id":  req.PackageID,
		"package_status":       "pending",
		"package_selected_at":  now,
		"package_activated_at": gorm.Expr("NULL"),
	}

	if err := config.DB.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save package selection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package selection saved", "package_id": req.PackageID})
}

func ListUsers(c *gin.Context) {
	var users []authmodels.User
	if err := config.DB.Preload("SelectedPackage").Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func ActivateUserPackage(c *gin.Context) {
	userID := c.Param("id")

	var user authmodels.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.SelectedPackageID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has not selected a package yet"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"package_status":       "active",
		"package_activated_at": now,
	}

	if err := config.DB.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package activated"})
}
