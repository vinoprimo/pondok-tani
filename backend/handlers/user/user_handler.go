package user

import (
	"net/http"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"

	"github.com/gin-gonic/gin"
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
