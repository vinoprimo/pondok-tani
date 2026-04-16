package sales

import (
	warehousemodels "pondok-tani-backend/models/warehouse"
	"time"
)

type SalesDetail struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	SalesOrderID     uint    `gorm:"not null;index" json:"sales_order_id"`
	WarehouseStockID uint    `gorm:"not null;index" json:"warehouse_stock_id"`
	Quantity         float64 `gorm:"type:numeric(12,2);not null" json:"quantity"`
	UnitPrice        float64 `gorm:"type:numeric(14,2);not null" json:"unit_price"`
	TotalPrice       float64 `gorm:"type:numeric(14,2);not null" json:"total_price"`

	SalesOrder     SalesOrder                     `gorm:"foreignKey:SalesOrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	WarehouseStock warehousemodels.WarehouseStock `gorm:"foreignKey:WarehouseStockID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SalesDetail) TableName() string {
	return "sales_details"
}
