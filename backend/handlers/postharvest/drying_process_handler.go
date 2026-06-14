package postharvest

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	harvestmodels "pondok-tani-backend/models/harvest"
	postharvestmodels "pondok-tani-backend/models/postharvest"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type completeDryingProcessInput struct {
	EndDate       string  `json:"end_date"`
	FinalQuantity float64 `json:"final_quantity"`
	Notes         *string `json:"notes"`
}

type dryingProcessResponse struct {
	ID              uint       `json:"id"`
	HarvestOutputID uint       `json:"harvest_output_id"`
	HarvestID       uint       `json:"harvest_id"`
	PlantBatchID    uint       `json:"plant_batch_id"`
	BatchCode       string     `json:"batch_code"`
	PackageName     string     `json:"package_name"`
	InitialQuantity float64    `json:"initial_quantity"`
	FinalQuantity   *float64   `json:"final_quantity,omitempty"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	Status          string     `json:"status"`
	Notes           *string    `json:"notes,omitempty"`
	HasGrading      bool       `json:"has_grading"`
}

func ListDryingProcesses(c *gin.Context) {
	var processes []postharvestmodels.DryingProcess
	query := config.DB.
		Preload("HarvestOutput").
		Preload("HarvestOutput.Harvest").
		Preload("HarvestOutput.Harvest.PlantBatch").
		Preload("HarvestOutput.Harvest.PlantBatch.Investment").
		Preload("HarvestOutput.Harvest.PlantBatch.Investment.Package").
		Order("start_date DESC")

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("drying_processes.status = ?", strings.ToLower(status))
	}

	if plantBatchValue := strings.TrimSpace(c.Query("plant_batch_id")); plantBatchValue != "" {
		batchID, err := strconv.ParseUint(plantBatchValue, 10, 64)
		if err != nil || batchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id"})
			return
		}
		query = query.
			Joins("JOIN harvest_outputs ON harvest_outputs.id = drying_processes.harvest_output_id").
			Joins("JOIN harvests ON harvests.id = harvest_outputs.harvest_id").
			Where("harvests.plant_batch_id = ?", uint(batchID))
	}

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN harvest_outputs ON harvest_outputs.id = drying_processes.harvest_output_id").
				Joins("JOIN harvests ON harvests.id = harvest_outputs.harvest_id").
				Joins("JOIN plant_batches ON plant_batches.id = harvests.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	if err := query.Find(&processes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drying processes"})
		return
	}

	gradingMap := map[uint]bool{}
	var gradingCounts []struct {
		DryingProcessID uint
		Total          int64
	}
	if err := config.DB.Table("grading_batches").
		Select("drying_process_id, COUNT(*) as total").
		Group("drying_process_id").
		Scan(&gradingCounts).Error; err == nil {
		for _, item := range gradingCounts {
			gradingMap[item.DryingProcessID] = item.Total > 0
		}
	}

	response := make([]dryingProcessResponse, 0, len(processes))
	for _, item := range processes {
		harvest := item.HarvestOutput.Harvest
		batch := harvest.PlantBatch
		response = append(response, dryingProcessResponse{
			ID:              item.ID,
			HarvestOutputID: item.HarvestOutputID,
			HarvestID:       harvest.ID,
			PlantBatchID:    harvest.PlantBatchID,
			BatchCode:       batch.BatchCode,
			PackageName:     batch.Investment.Package.PackageName,
			InitialQuantity: item.InitialQuantity,
			FinalQuantity:   item.FinalQuantity,
			StartDate:       item.StartDate,
			EndDate:         item.EndDate,
			Status:          item.Status,
			Notes:           item.Notes,
			HasGrading:      gradingMap[item.ID],
		})
	}

	c.JSON(http.StatusOK, response)
}

func CompleteDryingProcess(c *gin.Context) {
	idValue := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseUint(idValue, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drying process id"})
		return
	}

	var input completeDryingProcessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.FinalQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "final_quantity must be greater than 0"})
		return
	}

	var process postharvestmodels.DryingProcess
	if err := config.DB.Preload("HarvestOutput").First(&process, "id = ?", uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Drying process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load drying process"})
		return
	}

	if strings.ToLower(process.Status) != "ongoing" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Drying process already completed"})
		return
	}

	if input.FinalQuantity > process.InitialQuantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "final_quantity cannot exceed initial_quantity"})
		return
	}

	endDate := time.Now()
	if strings.TrimSpace(input.EndDate) != "" {
		if parsed, err := time.Parse(time.RFC3339, input.EndDate); err == nil {
			endDate = parsed
		} else if parsed, err := time.ParseInLocation("2006-01-02", input.EndDate, time.Local); err == nil {
			endDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be RFC3339 or YYYY-MM-DD"})
			return
		}
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		process.Status = "completed"
		process.EndDate = &endDate
		process.FinalQuantity = &input.FinalQuantity
		process.Notes = input.Notes

		if err := tx.Save(&process).Error; err != nil {
			return err
		}

		output := harvestmodels.HarvestOutput{}
		if err := tx.First(&output, "id = ?", process.HarvestOutputID).Error; err != nil {
			return err
		}
		output.Status = "ready_for_grading"
		if err := tx.Save(&output).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete drying process"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drying process completed"})
}
