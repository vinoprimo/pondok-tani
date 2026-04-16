package postharvest

import (
	harvestmodels "pondok-tani-backend/models/harvest"
	"time"
)

type DryingProcess struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	HarvestOutputID uint       `gorm:"not null;index" json:"harvest_output_id"`
	StartDate       time.Time  `gorm:"not null;index" json:"start_date"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	InitialQuantity float64    `gorm:"type:numeric(12,2);not null" json:"initial_quantity"`
	FinalQuantity   *float64   `gorm:"type:numeric(12,2)" json:"final_quantity,omitempty"`
	Status          string     `gorm:"type:varchar(30);not null;default:ongoing;index" json:"status"`
	Notes           *string    `gorm:"type:text" json:"notes,omitempty"`

	HarvestOutput harvestmodels.HarvestOutput `gorm:"foreignKey:HarvestOutputID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DryingProcess) TableName() string {
	return "drying_processes"
}
