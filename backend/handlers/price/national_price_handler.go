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
	PriceMinPerKg float64 `json:"price_min_per_kg"`
	PriceMaxPerKg float64 `json:"price_max_per_kg"`
	PricePerKg    float64 `json:"price_per_kg"`
	EffectiveDate string  `json:"effective_date" binding:"required"`
	HarvestType   string  `json:"harvest_type" binding:"required"`
}

func normalizePriceRange(req createOrUpdateNationalPriceRequest) (float64, float64, float64, bool) {
	priceMin := req.PriceMinPerKg
	priceMax := req.PriceMaxPerKg

	if priceMin <= 0 && priceMax <= 0 && req.PricePerKg > 0 {
		priceMin = req.PricePerKg
		priceMax = req.PricePerKg
	}

	if priceMin <= 0 || priceMax <= 0 || priceMin > priceMax {
		return 0, 0, 0, false
	}

	legacyPrice := (priceMin + priceMax) / 2
	return priceMin, priceMax, legacyPrice, true
}

func normalizeHarvestType(value string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "basah" || normalized == "kering" {
		return normalized, true
	}
	return "", false
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
	harvestType, isValidHarvestType := normalizeHarvestType(req.HarvestType)
	if province == "" || harvestType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "province and harvest_type are required"})
		return
	}
	if !isValidHarvestType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "harvest_type must be either basah or kering"})
		return
	}
	priceMinPerKg, priceMaxPerKg, legacyPricePerKg, validPrice := normalizePriceRange(req)
	if !validPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price_min_per_kg and price_max_per_kg must be greater than 0, and min cannot be greater than max"})
		return
	}

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid effective_date format. Use YYYY-MM-DD"})
		return
	}

	gradeID := req.GradeID
	if harvestType == "basah" {
		gradeID = nil
	}

	if gradeID != nil {
		var grade postharvestmodels.Grade
		if err := config.DB.First(&grade, "id = ?", *gradeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "grade_id is invalid"})
			return
		}
	}

	item := pricemodels.NationalPrice{
		GradeID:       gradeID,
		Province:      province,
		PriceMinPerKg: priceMinPerKg,
		PriceMaxPerKg: priceMaxPerKg,
		PricePerKg:    legacyPricePerKg,
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
	harvestType, isValidHarvestType := normalizeHarvestType(req.HarvestType)
	if province == "" || harvestType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "province and harvest_type are required"})
		return
	}
	if !isValidHarvestType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "harvest_type must be either basah or kering"})
		return
	}
	priceMinPerKg, priceMaxPerKg, legacyPricePerKg, validPrice := normalizePriceRange(req)
	if !validPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price_min_per_kg and price_max_per_kg must be greater than 0, and min cannot be greater than max"})
		return
	}

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid effective_date format. Use YYYY-MM-DD"})
		return
	}

	gradeID := req.GradeID
	if harvestType == "basah" {
		gradeID = nil
	}

	if gradeID != nil {
		var grade postharvestmodels.Grade
		if err := config.DB.First(&grade, "id = ?", *gradeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "grade_id is invalid"})
			return
		}
	}

	item.GradeID = gradeID
	item.Province = province
	item.PriceMinPerKg = priceMinPerKg
	item.PriceMaxPerKg = priceMaxPerKg
	item.PricePerKg = legacyPricePerKg
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
