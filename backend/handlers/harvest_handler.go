package handlers

import (
	"net/http"
	"pondok-tani-backend/config"
	"pondok-tani-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateHarvestOutput(c *gin.Context) {
	var input models.HarvestOutput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)

	c.JSON(http.StatusOK, input)
}