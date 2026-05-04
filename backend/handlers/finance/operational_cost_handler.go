package finance

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	coremodels "pondok-tani-backend/models/core"
	financemodels "pondok-tani-backend/models/finance"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createOrUpdateOperationalCostRequest struct {
	PlantBatchID uint    `json:"plant_batch_id" binding:"required,gt=0"`
	Category     string  `json:"category" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	CostDate     string  `json:"cost_date" binding:"required"`
	Note         string  `json:"note"`
}

type operationalCostResponse struct {
	ID           uint      `json:"id"`
	PlantBatchID uint      `json:"plant_batch_id"`
	BatchCode    string    `json:"batch_code"`
	PackageName  string    `json:"package_name"`
	Category     string    `json:"category"`
	Amount       float64   `json:"amount"`
	CostDate     time.Time `json:"cost_date"`
	Note         string    `json:"note"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func getUserRoleAndID(c *gin.Context) (string, string) {
	role := ""
	if roleValue, exists := c.Get("role"); exists {
		role, _ = roleValue.(string)
	}

	userID := ""
	if userIDValue, exists := c.Get("user_id"); exists {
		userID, _ = userIDValue.(string)
	}

	return role, userID
}

func applyPlantBatchOwnershipScope(c *gin.Context, query *gorm.DB) *gorm.DB {
	role, userID := getUserRoleAndID(c)
	if (role == "investor" || role == "mitra") && strings.TrimSpace(userID) != "" {
		return query.Joins("JOIN investments ON investments.id = plant_batches.investment_id").
			Where("investments.user_id = ?", userID)
	}
	return query
}

func applyOperationalCostOwnershipScope(c *gin.Context, query *gorm.DB) *gorm.DB {
	role, userID := getUserRoleAndID(c)
	if (role == "investor" || role == "mitra") && strings.TrimSpace(userID) != "" {
		return query.
			Joins("JOIN plant_batches ON plant_batches.id = operational_costs.plant_batch_id").
			Joins("JOIN investments ON investments.id = plant_batches.investment_id").
			Where("investments.user_id = ?", userID)
	}
	return query
}

func findAccessiblePlantBatch(c *gin.Context, plantBatchID uint) (coremodels.PlantBatch, error) {
	var plantBatch coremodels.PlantBatch
	query := config.DB.
		Preload("Investment").
		Preload("Investment.Package").
		Where("plant_batches.id = ?", plantBatchID)
	query = applyPlantBatchOwnershipScope(c, query)
	err := query.First(&plantBatch).Error
	return plantBatch, err
}

func mapOperationalCostResponse(item financemodels.OperationalCost) operationalCostResponse {
	createdBy := ""
	if item.CreatedBy != nil {
		createdBy = *item.CreatedBy
	}

	return operationalCostResponse{
		ID:           item.ID,
		PlantBatchID: item.PlantBatchID,
		BatchCode:    item.PlantBatch.BatchCode,
		PackageName:  item.PlantBatch.Investment.Package.PackageName,
		Category:     item.Category,
		Amount:       item.Amount,
		CostDate:     item.CostDate,
		Note:         item.Note,
		CreatedBy:    createdBy,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func ListOperationalCosts(c *gin.Context) {
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Model(&financemodels.OperationalCost{})

	query = applyOperationalCostOwnershipScope(c, query)

	if plantBatchIDQuery := strings.TrimSpace(c.Query("plant_batch_id")); plantBatchIDQuery != "" {
		plantBatchID, err := strconv.ParseUint(plantBatchIDQuery, 10, 64)
		if err != nil || plantBatchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id query"})
			return
		}
		query = query.Where("operational_costs.plant_batch_id = ?", uint(plantBatchID))
	}

	var items []financemodels.OperationalCost
	if err := query.Order("operational_costs.cost_date DESC, operational_costs.created_at DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch operational costs"})
		return
	}

	response := make([]operationalCostResponse, 0, len(items))
	for _, item := range items {
		response = append(response, mapOperationalCostResponse(item))
	}

	c.JSON(http.StatusOK, response)
}

func CreateOperationalCost(c *gin.Context) {
	var req createOrUpdateOperationalCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	category := strings.TrimSpace(req.Category)
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	costDate, err := time.Parse("2006-01-02", req.CostDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost_date format. Use YYYY-MM-DD"})
		return
	}

	plantBatch, err := findAccessiblePlantBatch(c, req.PlantBatchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant batch not found or inaccessible"})
		return
	}

	var createdBy *string
	if userIDValue, exists := c.Get("user_id"); exists {
		if userID, ok := userIDValue.(string); ok && strings.TrimSpace(userID) != "" {
			createdBy = &userID
		}
	}

	item := financemodels.OperationalCost{
		PlantBatchID: plantBatch.ID,
		Category:     category,
		Amount:       req.Amount,
		CostDate:     costDate,
		Note:         strings.TrimSpace(req.Note),
		CreatedBy:    createdBy,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create operational cost"})
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

	c.JSON(http.StatusCreated, mapOperationalCostResponse(item))
}

func UpdateOperationalCost(c *gin.Context) {
	query := config.DB.Model(&financemodels.OperationalCost{}).
		Where("operational_costs.id = ?", c.Param("id"))
	query = applyOperationalCostOwnershipScope(c, query)

	var item financemodels.OperationalCost
	if err := query.First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operational cost not found"})
		return
	}

	var req createOrUpdateOperationalCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	category := strings.TrimSpace(req.Category)
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	costDate, err := time.Parse("2006-01-02", req.CostDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost_date format. Use YYYY-MM-DD"})
		return
	}

	plantBatch, err := findAccessiblePlantBatch(c, req.PlantBatchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant batch not found or inaccessible"})
		return
	}

	item.PlantBatchID = plantBatch.ID
	item.Category = category
	item.Amount = req.Amount
	item.CostDate = costDate
	item.Note = strings.TrimSpace(req.Note)

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update operational cost"})
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

	c.JSON(http.StatusOK, mapOperationalCostResponse(item))
}

func DeleteOperationalCost(c *gin.Context) {
	query := config.DB.Model(&financemodels.OperationalCost{}).
		Where("operational_costs.id = ?", c.Param("id"))
	query = applyOperationalCostOwnershipScope(c, query)

	var item financemodels.OperationalCost
	if err := query.First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operational cost not found"})
		return
	}

	if err := config.DB.Delete(&financemodels.OperationalCost{}, "id = ?", item.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete operational cost"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operational cost deleted"})
}
