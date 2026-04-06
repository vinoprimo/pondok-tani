package handlers

import (
	"net/http"

	"pondok-tani-backend/config"
	"pondok-tani-backend/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User

	err := config.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}