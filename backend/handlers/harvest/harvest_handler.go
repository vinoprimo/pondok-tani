package harvest

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	harvestmodels "pondok-tani-backend/models/harvest"
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

type harvestListResponse struct {
	ID          uint      `json:"id"`
	PlantBatchID uint      `json:"plant_batch_id"`
	BatchCode   string    `json:"batch_code"`
	PackageName string    `json:"package_name"`
	HarvestDate time.Time `json:"harvest_date"`
	Notes       *string   `json:"notes"`
	CreatedBy   string    `json:"created_by"`
	CreatedByName string  `json:"created_by_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type harvestSummaryRow struct {
	Date          string  `json:"date"`
	TotalQuantity float64 `json:"total_quantity"`
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

func (h *HarvestHandler) ListHarvests(c *gin.Context) {
	var harvests []harvestmodels.Harvest
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Preload("Creator").
		Order("harvest_date DESC")

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN plant_batches ON plant_batches.id = harvests.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	if err := query.Find(&harvests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch harvest data"})
		return
	}

	response := make([]harvestListResponse, 0, len(harvests))
	for _, item := range harvests {
		createdByName := strings.TrimSpace(item.Creator.Name)
		response = append(response, harvestListResponse{
			ID:           item.ID,
			PlantBatchID: item.PlantBatchID,
			BatchCode:    item.PlantBatch.BatchCode,
			PackageName:  item.PlantBatch.Investment.Package.PackageName,
			HarvestDate:  item.HarvestDate,
			Notes:        item.Notes,
			CreatedBy:    item.CreatedBy,
			CreatedByName: createdByName,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *HarvestHandler) GetHarvestSummary(c *gin.Context) {
	var summary []harvestSummaryRow

	query := config.DB.Table("harvests").
		Select("DATE(harvests.harvest_date) as date, COALESCE(SUM(harvest_outputs.quantity), 0) as total_quantity").
		Joins("JOIN harvest_outputs ON harvest_outputs.harvest_id = harvests.id").
		Group("DATE(harvests.harvest_date)").
		Order("DATE(harvests.harvest_date) ASC")

	if batchIDValue := strings.TrimSpace(c.Query("plant_batch_id")); batchIDValue != "" {
		batchID, err := strconv.ParseUint(batchIDValue, 10, 64)
		if err != nil || batchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id"})
			return
		}
		query = query.Where("harvests.plant_batch_id = ?", uint(batchID))
	}

	if startDate := strings.TrimSpace(c.Query("start_date")); startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date"})
			return
		}
		query = query.Where("harvests.harvest_date >= ?", parsed)
	}

	if endDate := strings.TrimSpace(c.Query("end_date")); endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date"})
			return
		}
		query = query.Where("harvests.harvest_date <= ?", parsed)
	}

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN plant_batches ON plant_batches.id = harvests.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	if err := query.Scan(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch harvest summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
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
