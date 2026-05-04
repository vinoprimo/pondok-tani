package finance

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type OperationalCost struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PlantBatchID uint      `gorm:"not null;index" json:"plant_batch_id"`
	Category     string    `gorm:"type:varchar(80);not null;index" json:"category"`
	Amount       float64   `gorm:"type:numeric(14,2);not null" json:"amount"`
	CostDate     time.Time `gorm:"not null;index" json:"cost_date"`
	Note         string    `gorm:"type:text" json:"note"`
	CreatedBy    *string   `gorm:"type:uuid;index" json:"created_by,omitempty"`

	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Creator    *authmodels.User      `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (OperationalCost) TableName() string {
	return "operational_costs"
}
