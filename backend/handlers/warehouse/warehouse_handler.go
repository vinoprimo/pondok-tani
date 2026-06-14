package warehouse

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	warehousemodels "pondok-tani-backend/models/warehouse"

	"github.com/gin-gonic/gin"
)

type warehouseStockResponse struct {
	ID            uint      `json:"id"`
	PlantBatchID  uint      `json:"plant_batch_id"`
	BatchCode     string    `json:"batch_code"`
	PackageName   string    `json:"package_name"`
	UserID        string    `json:"user_id"`
	UserName      string    `json:"user_name"`
	GradeID       *uint     `json:"grade_id,omitempty"`
	GradeName     string    `json:"grade_name"`
	TotalQuantity float64   `json:"total_quantity"`
	Unit          string    `json:"unit"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type stockMovementResponse struct {
	ID           uint      `json:"id"`
	ReferenceType string   `json:"reference_type"`
	ReferenceID  uint      `json:"reference_id"`
	MovementType string    `json:"movement_type"`
	Quantity     float64   `json:"quantity"`
	MovementDate time.Time `json:"movement_date"`
	Notes        *string   `json:"notes,omitempty"`
	PlantBatchID uint      `json:"plant_batch_id"`
	BatchCode    string    `json:"batch_code"`
	PackageName  string    `json:"package_name"`
	GradeID      *uint     `json:"grade_id,omitempty"`
	GradeName    string    `json:"grade_name"`
}

type gradeBreakdownRow struct {
	GradeID   *uint   `json:"grade_id,omitempty"`
	GradeName string  `json:"grade_name"`
	Quantity  float64 `json:"quantity"`
}

type monthlyMovementRow struct {
	Month   string  `json:"month"`
	StockIn float64 `json:"stock_in"`
	StockOut float64 `json:"stock_out"`
}

type warehouseSummaryResponse struct {
	TotalQuantity   float64              `json:"total_quantity"`
	TotalBatches    int                  `json:"total_batches"`
	TotalGrades     int                  `json:"total_grades"`
	GradeBreakdown  []gradeBreakdownRow  `json:"grade_breakdown"`
	MonthlyMovements []monthlyMovementRow `json:"monthly_movements"`
}

func ListWarehouseStocks(c *gin.Context) {
	var stocks []warehousemodels.WarehouseStock
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Preload("Grade").
		Order("updated_at DESC")

	if plantBatchValue := strings.TrimSpace(c.Query("plant_batch_id")); plantBatchValue != "" {
		batchID, err := strconv.ParseUint(plantBatchValue, 10, 64)
		if err != nil || batchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id"})
			return
		}
		query = query.Where("warehouse_stocks.plant_batch_id = ?", uint(batchID))
	}

	if gradeValue := strings.TrimSpace(c.Query("grade_id")); gradeValue != "" {
		gradeID, err := strconv.ParseUint(gradeValue, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade_id"})
			return
		}
		if gradeID == 0 {
			query = query.Where("warehouse_stocks.grade_id IS NULL")
		} else {
			query = query.Where("warehouse_stocks.grade_id = ?", uint(gradeID))
		}
	}

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN plant_batches ON plant_batches.id = warehouse_stocks.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	if err := query.Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch warehouse stocks"})
		return
	}

	userIDs := make([]string, 0, len(stocks))
	userIDSet := make(map[string]struct{})
	for _, item := range stocks {
		userID := strings.TrimSpace(item.PlantBatch.Investment.UserID)
		if userID == "" {
			continue
		}
		if _, exists := userIDSet[userID]; exists {
			continue
		}
		userIDSet[userID] = struct{}{}
		userIDs = append(userIDs, userID)
	}

	userMap := map[string]string{}
	if len(userIDs) > 0 {
		var users []authmodels.User
		if err := config.DB.Select("id", "name").Where("id IN ?", userIDs).Find(&users).Error; err == nil {
			for _, user := range users {
				userMap[user.ID] = user.Name
			}
		}
	}

	response := make([]warehouseStockResponse, 0, len(stocks))
	for _, item := range stocks {
		gradeName := "Tidak diklasifikasi"
		if item.Grade != nil {
			gradeName = item.Grade.GradeName
		}
		userID := strings.TrimSpace(item.PlantBatch.Investment.UserID)
		userName := userMap[userID]
		response = append(response, warehouseStockResponse{
			ID:            item.ID,
			PlantBatchID:  item.PlantBatchID,
			BatchCode:     item.PlantBatch.BatchCode,
			PackageName:   item.PlantBatch.Investment.Package.PackageName,
			UserID:        userID,
			UserName:      userName,
			GradeID:       item.GradeID,
			GradeName:     gradeName,
			TotalQuantity: item.TotalQuantity,
			Unit:          item.Unit,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func ListStockMovements(c *gin.Context) {
	var movements []warehousemodels.StockMovement
	query := config.DB.
		Preload("WarehouseStock").
		Preload("WarehouseStock.PlantBatch").
		Preload("WarehouseStock.PlantBatch.Investment").
		Preload("WarehouseStock.PlantBatch.Investment.Package").
		Preload("WarehouseStock.Grade").
		Order("movement_date DESC")

	if plantBatchValue := strings.TrimSpace(c.Query("plant_batch_id")); plantBatchValue != "" {
		batchID, err := strconv.ParseUint(plantBatchValue, 10, 64)
		if err != nil || batchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id"})
			return
		}
		query = query.
			Joins("JOIN warehouse_stocks ON warehouse_stocks.id = stock_movements.warehouse_stock_id").
			Where("warehouse_stocks.plant_batch_id = ?", uint(batchID))
	}

	if gradeValue := strings.TrimSpace(c.Query("grade_id")); gradeValue != "" {
		gradeID, err := strconv.ParseUint(gradeValue, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade_id"})
			return
		}
		if gradeID == 0 {
			query = query.
				Joins("JOIN warehouse_stocks ON warehouse_stocks.id = stock_movements.warehouse_stock_id").
				Where("warehouse_stocks.grade_id IS NULL")
		} else {
			query = query.
				Joins("JOIN warehouse_stocks ON warehouse_stocks.id = stock_movements.warehouse_stock_id").
				Where("warehouse_stocks.grade_id = ?", uint(gradeID))
		}
	}

	if movementType := strings.TrimSpace(c.Query("movement_type")); movementType != "" {
		query = query.Where("stock_movements.movement_type = ?", strings.ToLower(movementType))
	}

	if startDate := strings.TrimSpace(c.Query("start_date")); startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date"})
			return
		}
		query = query.Where("stock_movements.movement_date >= ?", parsed)
	}

	if endDate := strings.TrimSpace(c.Query("end_date")); endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date"})
			return
		}
		query = query.Where("stock_movements.movement_date <= ?", parsed)
	}

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN warehouse_stocks ON warehouse_stocks.id = stock_movements.warehouse_stock_id").
				Joins("JOIN plant_batches ON plant_batches.id = warehouse_stocks.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	if err := query.Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock movements"})
		return
	}

	response := make([]stockMovementResponse, 0, len(movements))
	for _, item := range movements {
		gradeName := "-"
		if item.WarehouseStock.Grade != nil {
			gradeName = item.WarehouseStock.Grade.GradeName
		}
		response = append(response, stockMovementResponse{
			ID:            item.ID,
			ReferenceType: item.ReferenceType,
			ReferenceID:   item.ReferenceID,
			MovementType:  item.MovementType,
			Quantity:      item.Quantity,
			MovementDate:  item.MovementDate,
			Notes:         item.Notes,
			PlantBatchID:  item.WarehouseStock.PlantBatchID,
			BatchCode:     item.WarehouseStock.PlantBatch.BatchCode,
			PackageName:   item.WarehouseStock.PlantBatch.Investment.Package.PackageName,
			GradeID:       item.WarehouseStock.GradeID,
			GradeName:     gradeName,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetWarehouseSummary(c *gin.Context) {
	baseStocks := config.DB.Model(&warehousemodels.WarehouseStock{})
	baseMovements := config.DB.Model(&warehousemodels.StockMovement{})

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			baseStocks = baseStocks.
				Joins("JOIN plant_batches ON plant_batches.id = warehouse_stocks.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
			baseMovements = baseMovements.
				Joins("JOIN warehouse_stocks ON warehouse_stocks.id = stock_movements.warehouse_stock_id").
				Joins("JOIN plant_batches ON plant_batches.id = warehouse_stocks.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	var totalQuantity float64
	_ = baseStocks.Select("COALESCE(SUM(total_quantity), 0)").Scan(&totalQuantity).Error

	var totalBatches int
	_ = baseStocks.Select("COUNT(DISTINCT plant_batch_id)").Scan(&totalBatches).Error

	var totalGrades int
	_ = baseStocks.Select("COUNT(DISTINCT grade_id)").Scan(&totalGrades).Error

	var gradeBreakdown []gradeBreakdownRow
	gradeQuery := baseStocks.
		Select("warehouse_stocks.grade_id, COALESCE(grades.grade_name, 'Tidak diklasifikasi') as grade_name, COALESCE(SUM(warehouse_stocks.total_quantity), 0) as quantity").
		Joins("LEFT JOIN grades ON grades.id = warehouse_stocks.grade_id").
		Group("warehouse_stocks.grade_id, grades.grade_name")
	_ = gradeQuery.Scan(&gradeBreakdown).Error

	startMonth := time.Now().AddDate(0, -5, 0)
	startMonth = time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, time.Local)

	var movementRows []struct {
		MovementType string
		MovementDate time.Time
		Quantity     float64
	}
	_ = baseMovements.
		Select("movement_type, movement_date, quantity").
		Where("movement_date >= ?", startMonth).
		Find(&movementRows).Error

	monthMap := make(map[string]*monthlyMovementRow)
	for i := 0; i < 6; i++ {
		month := startMonth.AddDate(0, i, 0)
		key := month.Format("2006-01")
		monthMap[key] = &monthlyMovementRow{
			Month:   month.Format("Jan"),
			StockIn: 0,
			StockOut: 0,
		}
	}

	for _, row := range movementRows {
		key := row.MovementDate.Format("2006-01")
		item, exists := monthMap[key]
		if !exists {
			continue
		}
		if strings.ToLower(row.MovementType) == "in" {
			item.StockIn += row.Quantity
		} else {
			item.StockOut += row.Quantity
		}
	}

	monthlyMovements := make([]monthlyMovementRow, 0, len(monthMap))
	for i := 0; i < 6; i++ {
		month := startMonth.AddDate(0, i, 0)
		key := month.Format("2006-01")
		if item, ok := monthMap[key]; ok {
			monthlyMovements = append(monthlyMovements, *item)
		}
	}

	response := warehouseSummaryResponse{
		TotalQuantity:   totalQuantity,
		TotalBatches:    totalBatches,
		TotalGrades:     totalGrades,
		GradeBreakdown:  gradeBreakdown,
		MonthlyMovements: monthlyMovements,
	}

	c.JSON(http.StatusOK, response)
}
