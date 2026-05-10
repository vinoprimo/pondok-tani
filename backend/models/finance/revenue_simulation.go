package finance

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type RevenueSimulation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PlantBatchID uint      `gorm:"not null;index" json:"plant_batch_id"`
	Source       string    `gorm:"type:varchar(120);not null;index" json:"source"`
	Amount       float64   `gorm:"type:numeric(14,2);not null" json:"amount"`
	RevenueDate  time.Time `gorm:"not null;index" json:"revenue_date"`
	Note         string    `gorm:"type:text" json:"note"`
	Status       string    `gorm:"type:varchar(20);not null;default:simulasi;index" json:"status"`
	CreatedBy    *string   `gorm:"type:uuid;index" json:"created_by,omitempty"`

	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Creator    *authmodels.User      `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RevenueSimulation) TableName() string {
	return "revenue_simulations"
}
