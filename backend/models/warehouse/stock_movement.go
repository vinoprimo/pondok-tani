package warehouse

import "time"

type StockMovement struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	WarehouseStockID uint      `gorm:"not null;index" json:"warehouse_stock_id"`
	ReferenceType    string    `gorm:"type:varchar(50);not null;index" json:"reference_type"`
	ReferenceID      uint      `gorm:"not null;index" json:"reference_id"`
	MovementType     string    `gorm:"type:varchar(10);not null;index" json:"movement_type"`
	Quantity         float64   `gorm:"type:numeric(12,2);not null" json:"quantity"`
	MovementDate     time.Time `gorm:"not null;index" json:"movement_date"`
	Notes            *string   `gorm:"type:text" json:"notes,omitempty"`

	WarehouseStock WarehouseStock `gorm:"foreignKey:WarehouseStockID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (StockMovement) TableName() string {
	return "stock_movements"
}
