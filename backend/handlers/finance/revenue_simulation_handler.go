package finance

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	financemodels "pondok-tani-backend/models/finance"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createOrUpdateRevenueSimulationRequest struct {
	PlantBatchID uint    `json:"plant_batch_id" binding:"required,gt=0"`
	Source       string  `json:"source" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	RevenueDate  string  `json:"revenue_date" binding:"required"`
	Note         string  `json:"note"`
	Status       string  `json:"status"`
}

type revenueSimulationResponse struct {
	ID           uint      `json:"id"`
	PlantBatchID uint      `json:"plant_batch_id"`
	BatchCode    string    `json:"batch_code"`
	PackageName  string    `json:"package_name"`
	Source       string    `json:"source"`
	Amount       float64   `json:"amount"`
	RevenueDate  time.Time `json:"revenue_date"`
	Note         string    `json:"note"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func applyRevenueSimulationOwnershipScope(c *gin.Context, query *gorm.DB) *gorm.DB {
	role, userID := getUserRoleAndID(c)
	if (role == "investor" || role == "mitra") && strings.TrimSpace(userID) != "" {
		return query.
			Joins("JOIN plant_batches ON plant_batches.id = revenue_simulations.plant_batch_id").
			Joins("JOIN investments ON investments.id = plant_batches.investment_id").
			Where("investments.user_id = ?", userID)
	}
	return query
}

func mapRevenueSimulationResponse(item financemodels.RevenueSimulation) revenueSimulationResponse {
	createdBy := ""
	if item.CreatedBy != nil {
		createdBy = *item.CreatedBy
	}

	return revenueSimulationResponse{
		ID:           item.ID,
		PlantBatchID: item.PlantBatchID,
		BatchCode:    item.PlantBatch.BatchCode,
		PackageName:  item.PlantBatch.Investment.Package.PackageName,
		Source:       item.Source,
		Amount:       item.Amount,
		RevenueDate:  item.RevenueDate,
		Note:         item.Note,
		Status:       item.Status,
		CreatedBy:    createdBy,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func ListRevenueSimulations(c *gin.Context) {
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Model(&financemodels.RevenueSimulation{})

	query = applyRevenueSimulationOwnershipScope(c, query)

	if plantBatchIDQuery := strings.TrimSpace(c.Query("plant_batch_id")); plantBatchIDQuery != "" {
		plantBatchID, err := strconv.ParseUint(plantBatchIDQuery, 10, 64)
		if err != nil || plantBatchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id query"})
			return
		}
		query = query.Where("revenue_simulations.plant_batch_id = ?", uint(plantBatchID))
	}

	if statusQuery := strings.TrimSpace(c.Query("status")); statusQuery != "" {
		query = query.Where("revenue_simulations.status = ?", strings.ToLower(statusQuery))
	}

	var items []financemodels.RevenueSimulation
	if err := query.Order("revenue_simulations.revenue_date DESC, revenue_simulations.created_at DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch revenue simulations"})
		return
	}

	response := make([]revenueSimulationResponse, 0, len(items))
	for _, item := range items {
		response = append(response, mapRevenueSimulationResponse(item))
	}

	c.JSON(http.StatusOK, response)
}

func CreateRevenueSimulation(c *gin.Context) {
	var req createOrUpdateRevenueSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	source := strings.TrimSpace(req.Source)
	if source == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source is required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	revenueDate, err := time.Parse("2006-01-02", req.RevenueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revenue_date format. Use YYYY-MM-DD"})
		return
	}

	plantBatch, err := findAccessiblePlantBatch(c, req.PlantBatchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant batch not found or inaccessible"})
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "simulasi"
	}
	if status != "simulasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 'simulasi'"})
		return
	}

	var createdBy *string
	if userIDValue, exists := c.Get("user_id"); exists {
		if userID, ok := userIDValue.(string); ok && strings.TrimSpace(userID) != "" {
			createdBy = &userID
		}
	}

	item := financemodels.RevenueSimulation{
		PlantBatchID: plantBatch.ID,
		Source:       source,
		Amount:       req.Amount,
		RevenueDate:  revenueDate,
		Note:         strings.TrimSpace(req.Note),
		Status:       status,
		CreatedBy:    createdBy,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create revenue simulation"})
		return
	}

	if err := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		First(&item, "id = ?", item.ID).Error; err != nil {
		c.JSON(http.StatusCreated, item)
		return
	}

	c.JSON(http.StatusCreated, mapRevenueSimulationResponse(item))
}

func UpdateRevenueSimulation(c *gin.Context) {
	query := config.DB.Model(&financemodels.RevenueSimulation{}).
		Where("revenue_simulations.id = ?", c.Param("id"))
	query = applyRevenueSimulationOwnershipScope(c, query)

	var item financemodels.RevenueSimulation
	if err := query.First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Revenue simulation not found"})
		return
	}

	var req createOrUpdateRevenueSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	source := strings.TrimSpace(req.Source)
	if source == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source is required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	revenueDate, err := time.Parse("2006-01-02", req.RevenueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revenue_date format. Use YYYY-MM-DD"})
		return
	}

	plantBatch, err := findAccessiblePlantBatch(c, req.PlantBatchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant batch not found or inaccessible"})
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "simulasi"
	}
	if status != "simulasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 'simulasi'"})
		return
	}

	item.PlantBatchID = plantBatch.ID
	item.Source = source
	item.Amount = req.Amount
	item.RevenueDate = revenueDate
	item.Note = strings.TrimSpace(req.Note)
	item.Status = status

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update revenue simulation"})
		return
	}

	if err := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		First(&item, "id = ?", item.ID).Error; err != nil {
		c.JSON(http.StatusOK, item)
		return
	}

	c.JSON(http.StatusOK, mapRevenueSimulationResponse(item))
}

func DeleteRevenueSimulation(c *gin.Context) {
	query := config.DB.Model(&financemodels.RevenueSimulation{}).
		Where("revenue_simulations.id = ?", c.Param("id"))
	query = applyRevenueSimulationOwnershipScope(c, query)

	var item financemodels.RevenueSimulation
	if err := query.First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Revenue simulation not found"})
		return
	}

	if err := config.DB.Delete(&financemodels.RevenueSimulation{}, "id = ?", item.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete revenue simulation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Revenue simulation deleted"})
}
