package postharvest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	harvestmodels "pondok-tani-backend/models/harvest"
	postharvestmodels "pondok-tani-backend/models/postharvest"
	harvestrepo "pondok-tani-backend/repositories/harvest"
	harvestservice "pondok-tani-backend/services/harvest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type gradingDetailInput struct {
	GradeID  uint    `json:"grade_id"`
	Quantity float64 `json:"quantity"`
}

type createGradingBatchInput struct {
	DryingProcessID uint                `json:"drying_process_id"`
	GradingDate     string              `json:"grading_date"`
	Details         []gradingDetailInput `json:"details"`
}

type gradingDetailResponse struct {
	GradeID   uint    `json:"grade_id"`
	GradeName string  `json:"grade_name"`
	Quantity  float64 `json:"quantity"`
}

type gradingBatchResponse struct {
	ID              uint                   `json:"id"`
	DryingProcessID uint                   `json:"drying_process_id"`
	BatchCode       string                 `json:"batch_code"`
	GradingDate     time.Time              `json:"grading_date"`
	Status          string                 `json:"status"`
	PlantBatchID    uint                   `json:"plant_batch_id"`
	BatchCodePlant  string                 `json:"batch_code_plant"`
	PackageName     string                 `json:"package_name"`
	TotalQuantity   float64                `json:"total_quantity"`
	Details         []gradingDetailResponse `json:"details"`
}

func ListGradingBatches(c *gin.Context) {
	var batches []postharvestmodels.GradingBatch
	query := config.DB.
		Preload("DryingProcess").
		Preload("DryingProcess.HarvestOutput").
		Preload("DryingProcess.HarvestOutput.Harvest").
		Preload("DryingProcess.HarvestOutput.Harvest.PlantBatch").
		Preload("DryingProcess.HarvestOutput.Harvest.PlantBatch.Investment").
		Preload("DryingProcess.HarvestOutput.Harvest.PlantBatch.Investment.Package").
		Preload("DryingProcess.HarvestOutput.Harvest.PlantBatch").
		Preload("DryingProcess.HarvestOutput.Harvest.PlantBatch.Investment.Package").
		Order("grading_date DESC")

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN drying_processes ON drying_processes.id = grading_batches.drying_process_id").
				Joins("JOIN harvest_outputs ON harvest_outputs.id = drying_processes.harvest_output_id").
				Joins("JOIN harvests ON harvests.id = harvest_outputs.harvest_id").
				Joins("JOIN plant_batches ON plant_batches.id = harvests.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	if err := query.Find(&batches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grading batches"})
		return
	}

	batchIDs := make([]uint, 0, len(batches))
	for _, item := range batches {
		batchIDs = append(batchIDs, item.ID)
	}

	var details []postharvestmodels.GradingDetail
	if len(batchIDs) > 0 {
		if err := config.DB.Preload("Grade").Where("grading_batch_id IN ?", batchIDs).Find(&details).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grading details"})
			return
		}
	}

	detailsMap := make(map[uint][]gradingDetailResponse)
	for _, item := range details {
		detailsMap[item.GradingBatchID] = append(detailsMap[item.GradingBatchID], gradingDetailResponse{
			GradeID:   item.GradeID,
			GradeName: item.Grade.GradeName,
			Quantity:  item.Quantity,
		})
	}

	response := make([]gradingBatchResponse, 0, len(batches))
	for _, batch := range batches {
		process := batch.DryingProcess
		harvest := process.HarvestOutput.Harvest
		plantBatch := harvest.PlantBatch
		rows := detailsMap[batch.ID]
		totalQty := 0.0
		for _, row := range rows {
			totalQty += row.Quantity
		}

		response = append(response, gradingBatchResponse{
			ID:              batch.ID,
			DryingProcessID: batch.DryingProcessID,
			BatchCode:       batch.BatchCode,
			GradingDate:     batch.GradingDate,
			Status:          batch.Status,
			PlantBatchID:    harvest.PlantBatchID,
			BatchCodePlant:  plantBatch.BatchCode,
			PackageName:     plantBatch.Investment.Package.PackageName,
			TotalQuantity:   totalQty,
			Details:         rows,
		})
	}

	c.JSON(http.StatusOK, response)
}

func CreateGradingBatch(c *gin.Context) {
	var input createGradingBatchInput
	if err := c.ShouldBindJSON(&input); err != nil || input.DryingProcessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(input.Details) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "details are required"})
		return
	}

	var process postharvestmodels.DryingProcess
	if err := config.DB.
		Preload("HarvestOutput").
		Preload("HarvestOutput.Harvest").
		Preload("HarvestOutput.Harvest.PlantBatch").
		First(&process, "id = ?", input.DryingProcessID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Drying process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load drying process"})
		return
	}

	if strings.ToLower(process.Status) != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Drying process must be completed before grading"})
		return
	}

	var existing postharvestmodels.GradingBatch
	if err := config.DB.First(&existing, "drying_process_id = ?", process.ID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Grading batch already exists"})
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check grading batch"})
		return
	}

	if process.FinalQuantity == nil || *process.FinalQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Final quantity is required for grading"})
		return
	}

	total := 0.0
	for _, detail := range input.Details {
		if detail.GradeID == 0 || detail.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grading detail"})
			return
		}
		total += detail.Quantity
	}

	if total <= 0 || total-*process.FinalQuantity > 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total grading quantity cannot exceed final quantity"})
		return
	}

	gradingDate := time.Now()
	if strings.TrimSpace(input.GradingDate) != "" {
		if parsed, err := time.Parse(time.RFC3339, input.GradingDate); err == nil {
			gradingDate = parsed
		} else if parsed, err := time.ParseInLocation("2006-01-02", input.GradingDate, time.Local); err == nil {
			gradingDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "grading_date must be RFC3339 or YYYY-MM-DD"})
			return
		}
	}

	stockRepo := harvestrepo.NewGormWarehouseStockRepository(config.DB)
	movementRepo := harvestrepo.NewGormStockMovementRepository(config.DB)
	stockSvc := harvestservice.NewStockService(stockRepo, movementRepo)

	batchCode := fmt.Sprintf("GRD-%d-%s", process.ID, time.Now().Format("20060102150405"))
	plantBatchID := process.HarvestOutput.Harvest.PlantBatchID

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		gradingBatch := postharvestmodels.GradingBatch{
			DryingProcessID: process.ID,
			BatchCode:       batchCode,
			GradingDate:     gradingDate,
			Status:          "completed",
		}

		if err := tx.Create(&gradingBatch).Error; err != nil {
			return err
		}

		for _, detail := range input.Details {
			var grade postharvestmodels.Grade
			if err := tx.First(&grade, "id = ?", detail.GradeID).Error; err != nil {
				return err
			}

			entity := postharvestmodels.GradingDetail{
				GradingBatchID: gradingBatch.ID,
				GradeID:        detail.GradeID,
				Quantity:       detail.Quantity,
			}

			if err := tx.Create(&entity).Error; err != nil {
				return err
			}

			if err := stockSvc.RecordIncomingFromGradingDetail(c.Request.Context(), tx, entity.ID, plantBatchID, detail.GradeID, detail.Quantity); err != nil {
				return err
			}
		}

		output := harvestmodels.HarvestOutput{}
		if err := tx.First(&output, "id = ?", process.HarvestOutputID).Error; err != nil {
			return err
		}
		output.Status = "graded"
		if err := tx.Save(&output).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create grading batch"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Grading batch created"})
}
