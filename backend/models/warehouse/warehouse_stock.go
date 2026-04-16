package warehouse

import (
	coremodels "pondok-tani-backend/models/core"
	postharvestmodels "pondok-tani-backend/models/postharvest"
	"time"
)

type WarehouseStock struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	PlantBatchID  uint    `gorm:"not null;index" json:"plant_batch_id"`
	GradeID       *uint   `gorm:"index" json:"grade_id,omitempty"`
	TotalQuantity float64 `gorm:"type:numeric(12,2);not null;default:0" json:"total_quantity"`
	Unit          string  `gorm:"type:varchar(10);not null;default:kg" json:"unit"`

	PlantBatch coremodels.PlantBatch    `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Grade      *postharvestmodels.Grade `gorm:"foreignKey:GradeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (WarehouseStock) TableName() string {
	return "warehouse_stocks"
}
