package harvest

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	harvestmodels "pondok-tani-backend/models/harvest"
	maintenancemodels "pondok-tani-backend/models/maintenance"
	postharvestmodels "pondok-tani-backend/models/postharvest"
	harvestrepo "pondok-tani-backend/repositories/harvest"
	harvestservice "pondok-tani-backend/services/harvest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HarvestRequestHandler struct{}

type createHarvestRequestInput struct {
	PlantMonitoringID uint `json:"plant_monitoring_id"`
}

type validateHarvestRequestInput struct {
	HarvestDate   string  `json:"harvest_date"`
	Notes         *string `json:"notes"`
	TotalQuantity float64 `json:"total_quantity"`
	WetQuantity   float64 `json:"wet_quantity"`
	DryQuantity   float64 `json:"dry_quantity"`
}

type harvestRequestResponse struct {
	ID               uint       `json:"id"`
	PlantMonitoringID uint       `json:"plant_monitoring_id"`
	PlantBatchID     uint       `json:"plant_batch_id"`
	BatchCode        string     `json:"batch_code"`
	PackageName      string     `json:"package_name"`
	RequestedBy      string     `json:"requested_by"`
	RequestedByName  string     `json:"requested_by_name"`
	Status           string     `json:"status"`
	Phase            string     `json:"phase"`
	Progress         int        `json:"progress"`
	HealthStatus     string     `json:"health_status"`
	Disease          string     `json:"disease"`
	DiseaseNote      string     `json:"disease_note"`
	AffectedCount    int        `json:"affected_count"`
	TotalPlants      int        `json:"total_plants"`
	Note             string     `json:"note"`
	PhotoURL         string     `json:"photo_url"`
	MonitoringDate   time.Time  `json:"monitoring_date"`
	TotalQuantity    float64    `json:"total_quantity"`
	WetQuantity      float64    `json:"wet_quantity"`
	DryQuantity      float64    `json:"dry_quantity"`
	HarvestID        *uint      `json:"harvest_id,omitempty"`
	HarvestDate      *time.Time `json:"harvest_date,omitempty"`
	ValidatedBy      string     `json:"validated_by,omitempty"`
	ValidatedByName  string     `json:"validated_by_name,omitempty"`
	ValidatedAt      *time.Time `json:"validated_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func NewHarvestRequestHandler() *HarvestRequestHandler {
	return &HarvestRequestHandler{}
}

func (h *HarvestRequestHandler) CreateHarvestRequest(c *gin.Context) {
	var input createHarvestRequestInput
	if err := c.ShouldBindJSON(&input); err != nil || input.PlantMonitoringID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plant_monitoring_id is required"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	userID, ok := userIDValue.(string)
	if !exists || !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	roleValue, _ := c.Get("role")
	role, _ := roleValue.(string)

	var monitoring maintenancemodels.PlantMonitoring
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package")
	if role == "investor" || role == "mitra" {
		query = query.
			Joins("JOIN plant_batches ON plant_batches.id = plant_monitorings.plant_batch_id").
			Joins("JOIN investments ON investments.id = plant_batches.investment_id").
			Where("investments.user_id = ?", userID)
	}

	if err := query.First(&monitoring, "plant_monitorings.id = ?", input.PlantMonitoringID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plant monitoring not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load plant monitoring"})
		return
	}

	if strings.ToLower(strings.TrimSpace(monitoring.Phase)) != "panen" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plant monitoring phase must be panen"})
		return
	}

	var existing harvestmodels.HarvestRequest
	if err := config.DB.
		Where("plant_batch_id = ? AND status = ?", monitoring.PlantBatchID, "pending").
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Harvest request already pending for this batch"})
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing harvest requests"})
		return
	}

	entity := harvestmodels.HarvestRequest{
		PlantMonitoringID: monitoring.ID,
		PlantBatchID:      monitoring.PlantBatchID,
		RequestedBy:       userID,
		Status:            "pending",
		Phase:             monitoring.Phase,
		Progress:          monitoring.Progress,
		HealthStatus:      monitoring.HealthStatus,
		Disease:           monitoring.Disease,
		DiseaseNote:       monitoring.DiseaseNote,
		AffectedCount:     monitoring.AffectedCount,
		TotalPlants:       monitoring.TotalPlants,
		Note:              monitoring.Note,
		PhotoURL:          monitoring.PhotoURL,
		MonitoringDate:    monitoring.CreatedAt,
	}

	if err := config.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create harvest request"})
		return
	}

	entity.PlantBatch = monitoring.PlantBatch
	c.JSON(http.StatusCreated, buildHarvestRequestResponse(entity))
}

func (h *HarvestRequestHandler) ListHarvestRequests(c *gin.Context) {
	var requests []harvestmodels.HarvestRequest
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Preload("Requester").
		Preload("Validator").
		Preload("Harvest")

	if statusQuery := strings.TrimSpace(c.Query("status")); statusQuery != "" {
		query = query.Where("status = ?", strings.ToLower(statusQuery))
	}

	roleValue, _ := c.Get("role")
	role, _ := roleValue.(string)
	if role == "investor" || role == "mitra" {
		userIDValue, _ := c.Get("user_id")
		userID, _ := userIDValue.(string)
		query = query.Where("requested_by = ?", userID)
	}

	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch harvest requests"})
		return
	}

	response := make([]harvestRequestResponse, 0, len(requests))
	for _, item := range requests {
		response = append(response, buildHarvestRequestResponse(item))
	}

	c.JSON(http.StatusOK, response)
}

func (h *HarvestRequestHandler) ValidateHarvestRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || requestID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harvest request id"})
		return
	}

	var input validateHarvestRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	harvestDate, err := parseHarvestDateInput(input.HarvestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "harvest_date must be RFC3339 or YYYY-MM-DD"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	userID, ok := userIDValue.(string)
	if !exists || !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	var request harvestmodels.HarvestRequest
	if err := config.DB.First(&request, "id = ?", uint(requestID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Harvest request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load harvest request"})
		return
	}

	if request.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harvest request already processed"})
		return
	}

	totalQty := input.TotalQuantity
	wetQty := input.WetQuantity
	dryQty := input.DryQuantity

	if wetQty < 0 || dryQty < 0 || totalQty < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity cannot be negative"})
		return
	}

	combinedQty := wetQty + dryQty
	if totalQty == 0 && combinedQty > 0 {
		totalQty = combinedQty
	}

	if totalQty <= 0 || combinedQty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total quantity and output quantity are required"})
		return
	}

	if math.Abs(totalQty-combinedQty) > 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total quantity must equal wet + dry quantity"})
		return
	}

	now := time.Now()
	stockRepo := harvestrepo.NewGormWarehouseStockRepository(config.DB)
	movementRepo := harvestrepo.NewGormStockMovementRepository(config.DB)
	stockSvc := harvestservice.NewStockService(stockRepo, movementRepo)

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		harvest := harvestmodels.Harvest{
			PlantBatchID: request.PlantBatchID,
			HarvestDate:  harvestDate,
			CreatedBy:    userID,
			Notes:        input.Notes,
		}

		if err := tx.Create(&harvest).Error; err != nil {
			return err
		}

		if wetQty > 0 {
			output := harvestmodels.HarvestOutput{
				HarvestID: harvest.ID,
				Jenis:     "basah",
				Quantity:  wetQty,
				Status:    "stocked",
			}
			if err := tx.Create(&output).Error; err != nil {
				return err
			}
			if err := stockSvc.RecordIncomingFromHarvestOutput(c.Request.Context(), tx, output.ID, harvest.PlantBatchID, nil, wetQty); err != nil {
				return err
			}
		}

		if dryQty > 0 {
			output := harvestmodels.HarvestOutput{
				HarvestID: harvest.ID,
				Jenis:     "kering",
				Quantity:  dryQty,
				Status:    "pending_drying",
			}
			if err := tx.Create(&output).Error; err != nil {
				return err
			}

			drying := postharvestmodels.DryingProcess{
				HarvestOutputID: output.ID,
				StartDate:       now,
				InitialQuantity: dryQty,
				Status:          "ongoing",
			}
			if err := tx.Create(&drying).Error; err != nil {
				return err
			}
		}

		request.Status = "validated"
		request.ValidatedBy = &userID
		request.ValidatedAt = &now
		request.HarvestID = &harvest.ID
		request.TotalQuantity = totalQty
		request.WetQuantity = wetQty
		request.DryQuantity = dryQty

		if err := tx.Save(&request).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate harvest request"})
		return
	}

	var refreshed harvestmodels.HarvestRequest
	if err := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Preload("Requester").
		Preload("Validator").
		Preload("Harvest").
		First(&refreshed, "id = ?", request.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Harvest request validated", "data": buildHarvestRequestResponse(request)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Harvest request validated", "data": buildHarvestRequestResponse(refreshed)})
}

func parseHarvestDateInput(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, errors.New("harvest_date is required")
	}

	if parsed, err := time.Parse(time.RFC3339, trimmed); err == nil {
		return parsed, nil
	}

	if parsed, err := time.ParseInLocation("2006-01-02", trimmed, time.Local); err == nil {
		return parsed, nil
	}

	return time.Time{}, errors.New("invalid harvest_date")
}

func buildHarvestRequestResponse(item harvestmodels.HarvestRequest) harvestRequestResponse {
	requestedByName := strings.TrimSpace(item.Requester.Name)
	validatedByName := strings.TrimSpace(item.Validator.Name)

	var harvestDate *time.Time
	if item.HarvestID != nil && !item.Harvest.HarvestDate.IsZero() {
		date := item.Harvest.HarvestDate
		harvestDate = &date
	}

	validatedBy := ""
	if item.ValidatedBy != nil {
		validatedBy = *item.ValidatedBy
	}

	return harvestRequestResponse{
		ID:               item.ID,
		PlantMonitoringID: item.PlantMonitoringID,
		PlantBatchID:     item.PlantBatchID,
		BatchCode:        item.PlantBatch.BatchCode,
		PackageName:      item.PlantBatch.Investment.Package.PackageName,
		RequestedBy:      item.RequestedBy,
		RequestedByName:  requestedByName,
		Status:           item.Status,
		Phase:            item.Phase,
		Progress:         item.Progress,
		HealthStatus:     item.HealthStatus,
		Disease:          item.Disease,
		DiseaseNote:      item.DiseaseNote,
		AffectedCount:    item.AffectedCount,
		TotalPlants:      item.TotalPlants,
		Note:             item.Note,
		PhotoURL:         item.PhotoURL,
		MonitoringDate:   item.MonitoringDate,
		TotalQuantity:    item.TotalQuantity,
		WetQuantity:      item.WetQuantity,
		DryQuantity:      item.DryQuantity,
		HarvestID:        item.HarvestID,
		HarvestDate:      harvestDate,
		ValidatedBy:      validatedBy,
		ValidatedByName:  validatedByName,
		ValidatedAt:      item.ValidatedAt,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}
