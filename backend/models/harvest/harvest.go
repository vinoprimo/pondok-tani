package harvest

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type Harvest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PlantBatchID uint      `gorm:"not null;index" json:"plant_batch_id"`
	HarvestDate  time.Time `gorm:"not null;index" json:"harvest_date"`
	CreatedBy    string    `gorm:"type:uuid;not null;index" json:"created_by"`
	Notes        *string   `gorm:"type:text" json:"notes,omitempty"`

	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Creator    authmodels.User       `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Outputs []HarvestOutput `gorm:"foreignKey:HarvestID" json:"outputs,omitempty"`
}

func (Harvest) TableName() string {
	return "harvests"
}
