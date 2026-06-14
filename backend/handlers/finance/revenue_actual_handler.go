package finance

import (
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	salesmodels "pondok-tani-backend/models/sales"

	"github.com/gin-gonic/gin"
)

type revenueActualResponse struct {
	ID          uint      `json:"id"`
	PlantBatchID uint     `json:"plant_batch_id"`
	BatchCode   string    `json:"batch_code"`
	PackageName string    `json:"package_name"`
	Source      string    `json:"source"`
	Amount      float64   `json:"amount"`
	RevenueDate time.Time `json:"revenue_date"`
	Note        string    `json:"note"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func ListActualRevenues(c *gin.Context) {
	query := config.DB.
		Preload("SalesOrder").
		Preload("WarehouseStock").
		Preload("WarehouseStock.PlantBatch").
		Preload("WarehouseStock.PlantBatch.Investment").
		Preload("WarehouseStock.PlantBatch.Investment.Package").
		Joins("JOIN sales_orders ON sales_orders.id = sales_details.sales_order_id")

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("sales_orders.status = ?", strings.ToLower(status))
	}

	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			userIDValue, _ := c.Get("user_id")
			userID, _ := userIDValue.(string)
			query = query.
				Joins("JOIN warehouse_stocks ON warehouse_stocks.id = sales_details.warehouse_stock_id").
				Joins("JOIN plant_batches ON plant_batches.id = warehouse_stocks.plant_batch_id").
				Joins("JOIN investments ON investments.id = plant_batches.investment_id").
				Where("investments.user_id = ?", userID)
		}
	}

	var details []salesmodels.SalesDetail
	if err := query.Order("sales_orders.sales_date DESC, sales_details.id DESC").Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch revenues"})
		return
	}

	response := make([]revenueActualResponse, 0, len(details))
	for _, detail := range details {
		batch := detail.WarehouseStock.PlantBatch
		packageName := ""
		if batch.Investment.Package.PackageName != "" {
			packageName = batch.Investment.Package.PackageName
		}

		response = append(response, revenueActualResponse{
			ID:           detail.ID,
			PlantBatchID: batch.ID,
			BatchCode:    batch.BatchCode,
			PackageName:  packageName,
			Source:       "Penjualan panen",
			Amount:       detail.TotalPrice,
			RevenueDate:  detail.SalesOrder.SalesDate,
			Note:         detail.SalesOrder.Buyer,
			Status:       "realisasi",
			CreatedAt:    detail.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
