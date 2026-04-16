package harvest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	harvestrepo "pondok-tani-backend/repositories/harvest"
	harvestservice "pondok-tani-backend/services/harvest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HarvestHandler struct {
	harvestService harvestservice.HarvestService
	outputService  harvestservice.HarvestOutputService
}

type createHarvestRequest struct {
	PlantBatchID uint    `json:"plant_batch_id"`
	HarvestDate  string  `json:"harvest_date"`
	Notes        *string `json:"notes"`
}

type createHarvestOutputRequest struct {
	Jenis    string  `json:"jenis"`
	Quantity float64 `json:"quantity"`
}

func NewHarvestHandler(db *gorm.DB) *HarvestHandler {
	harvestRepo := harvestrepo.NewGormHarvestRepository(db)
	outputRepo := harvestrepo.NewGormHarvestOutputRepository(db)
	stockRepo := harvestrepo.NewGormWarehouseStockRepository(db)
	movementRepo := harvestrepo.NewGormStockMovementRepository(db)

	stockSvc := harvestservice.NewStockService(stockRepo, movementRepo)

	return &HarvestHandler{
		harvestService: harvestservice.NewHarvestService(harvestRepo),
		outputService:  harvestservice.NewHarvestOutputService(db, harvestRepo, outputRepo, stockSvc),
	}
}

func (h *HarvestHandler) CreateHarvest(c *gin.Context) {
	var req createHarvestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdBy, _ := c.Get("user_id")
	createdByStr, _ := createdBy.(string)

	harvestDate, err := time.Parse(time.RFC3339, req.HarvestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "harvest_date must be RFC3339 format"})
		return
	}

	result, err := h.harvestService.CreateHarvest(c.Request.Context(), harvestservice.CreateHarvestInput{
		PlantBatchID: req.PlantBatchID,
		HarvestDate:  harvestDate,
		CreatedBy:    createdByStr,
		Notes:        req.Notes,
	})
	if err != nil {
		handleHarvestServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *HarvestHandler) CreateHarvestOutput(c *gin.Context) {
	harvestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harvest id"})
		return
	}

	var req createHarvestOutputRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.outputService.CreateHarvestOutput(c.Request.Context(), uint(harvestID), harvestservice.CreateHarvestOutputInput{
		Jenis:    req.Jenis,
		Quantity: req.Quantity,
	})
	if err != nil {
		handleHarvestServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func handleHarvestServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, harvestservice.ErrInvalidHarvestDate),
		errors.Is(err, harvestservice.ErrInvalidBatchID),
		errors.Is(err, harvestservice.ErrInvalidJenis),
		errors.Is(err, harvestservice.ErrInvalidBerat):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, harvestservice.ErrHarvestNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
