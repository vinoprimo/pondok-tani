package price

import (
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	postharvestmodels "pondok-tani-backend/models/postharvest"
	pricemodels "pondok-tani-backend/models/price"

	"github.com/gin-gonic/gin"
)

type createOrUpdateNationalPriceRequest struct {
	GradeID       *uint   `json:"grade_id"`
	Province      string  `json:"province" binding:"required"`
	PricePerKg    float64 `json:"price_per_kg" binding:"required"`
	EffectiveDate string  `json:"effective_date" binding:"required"`
	HarvestType   string  `json:"harvest_type" binding:"required"`
}

func ListNationalPrices(c *gin.Context) {
	var items []pricemodels.NationalPrice
	if err := config.DB.
		Preload("Grade").
		Order("effective_date DESC, national_price_id DESC").
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch national prices"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func CreateNationalPrice(c *gin.Context) {
	var req createOrUpdateNationalPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	province := strings.TrimSpace(req.Province)
	harvestType := strings.TrimSpace(req.HarvestType)
	if province == "" || harvestType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "province and harvest_type are required"})
		return
	}
	if req.PricePerKg <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price_per_kg must be greater than 0"})
		return
	}

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid effective_date format. Use YYYY-MM-DD"})
		return
	}

	if req.GradeID != nil {
		var grade postharvestmodels.Grade
		if err := config.DB.First(&grade, "id = ?", *req.GradeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "grade_id is invalid"})
			return
		}
	}

	item := pricemodels.NationalPrice{
		GradeID:       req.GradeID,
		Province:      province,
		PricePerKg:    req.PricePerKg,
		EffectiveDate: effectiveDate,
		HarvestType:   harvestType,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create national price"})
		return
	}

	if err := config.DB.Preload("Grade").First(&item, "national_price_id = ?", item.NationalPriceID).Error; err != nil {
		c.JSON(http.StatusCreated, item)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func UpdateNationalPrice(c *gin.Context) {
	var item pricemodels.NationalPrice
	if err := config.DB.First(&item, "national_price_id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "National price not found"})
		return
	}

	var req createOrUpdateNationalPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	province := strings.TrimSpace(req.Province)
	harvestType := strings.TrimSpace(req.HarvestType)
	if province == "" || harvestType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "province and harvest_type are required"})
		return
	}
	if req.PricePerKg <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price_per_kg must be greater than 0"})
		return
	}

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid effective_date format. Use YYYY-MM-DD"})
		return
	}

	if req.GradeID != nil {
		var grade postharvestmodels.Grade
		if err := config.DB.First(&grade, "id = ?", *req.GradeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "grade_id is invalid"})
			return
		}
	}

	item.GradeID = req.GradeID
	item.Province = province
	item.PricePerKg = req.PricePerKg
	item.EffectiveDate = effectiveDate
	item.HarvestType = harvestType

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update national price"})
		return
	}

	if err := config.DB.Preload("Grade").First(&item, "national_price_id = ?", item.NationalPriceID).Error; err != nil {
		c.JSON(http.StatusOK, item)
		return
	}

	c.JSON(http.StatusOK, item)
}

func DeleteNationalPrice(c *gin.Context) {
	if err := config.DB.Delete(&pricemodels.NationalPrice{}, "national_price_id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete national price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "National price deleted"})
}
