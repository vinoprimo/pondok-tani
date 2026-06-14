package sales

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	salesmodels "pondok-tani-backend/models/sales"
	warehousemodels "pondok-tani-backend/models/warehouse"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createSalesDetailInput struct {
	WarehouseStockID uint    `json:"warehouse_stock_id"`
	Quantity         float64 `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
}

type createSalesOrderInput struct {
	SalesDate string                 `json:"sales_date"`
	Buyer     string                 `json:"buyer"`
	Details   []createSalesDetailInput `json:"details"`
}

type salesDetailResponse struct {
	ID               uint    `json:"id"`
	WarehouseStockID uint    `json:"warehouse_stock_id"`
	PlantBatchID     uint    `json:"plant_batch_id"`
	BatchCode        string  `json:"batch_code"`
	PackageName      string  `json:"package_name"`
	GradeName        string  `json:"grade_name"`
	Quantity         float64 `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
	TotalPrice       float64 `json:"total_price"`
}

type salesOrderResponse struct {
	ID          uint                 `json:"id"`
	OrderNumber string               `json:"order_number"`
	SalesDate   time.Time            `json:"sales_date"`
	Buyer       string               `json:"buyer"`
	TotalAmount float64              `json:"total_amount"`
	Status      string               `json:"status"`
	CreatedBy   string               `json:"created_by"`
	CreatedAt   time.Time            `json:"created_at"`
	Details     []salesDetailResponse `json:"details"`
}

func ListSalesOrders(c *gin.Context) {
	query := config.DB.
		Preload("SalesOrder").
		Preload("WarehouseStock").
		Preload("WarehouseStock.PlantBatch").
		Preload("WarehouseStock.PlantBatch.Investment").
		Preload("WarehouseStock.PlantBatch.Investment.Package").
		Preload("WarehouseStock.Grade").
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
	if err := query.Order("sales_orders.sales_date DESC, sales_details.id ASC").Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales orders"})
		return
	}

	orderMap := make(map[uint]*salesOrderResponse)
	for _, detail := range details {
		order := detail.SalesOrder
		orderRow, exists := orderMap[order.ID]
		if !exists {
			orderRow = &salesOrderResponse{
				ID:          order.ID,
				OrderNumber: order.OrderNumber,
				SalesDate:   order.SalesDate,
				Buyer:       order.Buyer,
				TotalAmount: order.TotalAmount,
				Status:      order.Status,
				CreatedBy:   order.CreatedBy,
				CreatedAt:   order.CreatedAt,
				Details:     []salesDetailResponse{},
			}
			orderMap[order.ID] = orderRow
		}

		gradeName := "Tidak diklasifikasi"
		if detail.WarehouseStock.Grade != nil {
			gradeName = detail.WarehouseStock.Grade.GradeName
		}

		orderRow.Details = append(orderRow.Details, salesDetailResponse{
			ID:               detail.ID,
			WarehouseStockID: detail.WarehouseStockID,
			PlantBatchID:     detail.WarehouseStock.PlantBatchID,
			BatchCode:        detail.WarehouseStock.PlantBatch.BatchCode,
			PackageName:      detail.WarehouseStock.PlantBatch.Investment.Package.PackageName,
			GradeName:        gradeName,
			Quantity:         detail.Quantity,
			UnitPrice:        detail.UnitPrice,
			TotalPrice:       detail.TotalPrice,
		})
	}

	response := make([]salesOrderResponse, 0, len(orderMap))
	for _, order := range orderMap {
		response = append(response, *order)
	}

	c.JSON(http.StatusOK, response)
}

func CreateSalesOrder(c *gin.Context) {
	var input createSalesOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	buyer := strings.TrimSpace(input.Buyer)
	if buyer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "buyer is required"})
		return
	}
	if len(input.Details) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "details are required"})
		return
	}

	salesDate, err := time.Parse("2006-01-02", input.SalesDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sales_date format. Use YYYY-MM-DD"})
		return
	}

	createdByValue, exists := c.Get("user_id")
	createdBy, ok := createdByValue.(string)
	if !exists || !ok || strings.TrimSpace(createdBy) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	totalAmount := 0.0
	for _, detail := range input.Details {
		if detail.WarehouseStockID == 0 || detail.Quantity <= 0 || detail.UnitPrice <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sales detail"})
			return
		}
		totalAmount += detail.Quantity * detail.UnitPrice
	}

	orderNumber := fmt.Sprintf("SO-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%100000)

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		order := salesmodels.SalesOrder{
			OrderNumber: orderNumber,
			SalesDate:   salesDate,
			Buyer:       buyer,
			TotalAmount: totalAmount,
			Status:      "completed",
			CreatedBy:   createdBy,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, detail := range input.Details {
			var stock warehousemodels.WarehouseStock
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&stock, "id = ?", detail.WarehouseStockID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("stock not found")
				}
				return err
			}

			if stock.TotalQuantity < detail.Quantity {
				return errors.New("insufficient stock")
			}

			stock.TotalQuantity -= detail.Quantity
			if err := tx.Save(&stock).Error; err != nil {
				return err
			}

			totalPrice := detail.Quantity * detail.UnitPrice
			entity := salesmodels.SalesDetail{
				SalesOrderID:     order.ID,
				WarehouseStockID: detail.WarehouseStockID,
				Quantity:         detail.Quantity,
				UnitPrice:        detail.UnitPrice,
				TotalPrice:       totalPrice,
			}

			if err := tx.Create(&entity).Error; err != nil {
				return err
			}

			movement := warehousemodels.StockMovement{
				WarehouseStockID: stock.ID,
				ReferenceType:    "sales_detail",
				ReferenceID:      entity.ID,
				MovementType:     "out",
				Quantity:         detail.Quantity,
				MovementDate:     time.Now(),
			}

			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		if err.Error() == "stock not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse stock not found"})
			return
		}
		if err.Error() == "insufficient stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sales order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Sales order created"})
}
