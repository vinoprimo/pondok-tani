package harvest

import (
	"net/http"

	"pondok-tani-backend/config"
	harvestmodels "pondok-tani-backend/models/harvest"

	"github.com/gin-gonic/gin"
)

func CreateHarvestOutput(c *gin.Context) {
	var input harvestmodels.HarvestOutput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)

	c.JSON(http.StatusOK, input)
}
