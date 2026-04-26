package maintenance

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type PlantMonitoring struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	PlantBatchID uint `gorm:"not null;index" json:"plant_batch_id"`

	// Growth
	Phase    string `gorm:"type:varchar(50);not null" json:"phase"`
	Progress int    `gorm:"not null" json:"progress"`

	// Health
	HealthStatus string `gorm:"type:varchar(30);not null" json:"health_status"`
	Disease      string `gorm:"type:varchar(50)" json:"disease"`

	// Additional notes for custom disease details
	DiseaseNote string `gorm:"type:text" json:"disease_note"`

	// Aggregation
	AffectedCount int `json:"affected_count"`
	TotalPlants   int `json:"total_plants"`

	// Info
	Note string `gorm:"type:text" json:"note"`

	CreatedBy *string `gorm:"type:uuid;index" json:"created_by,omitempty"`
	PhotoURL  string  `gorm:"type:text" json:"photo_url"`

	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Creator    authmodels.User       `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PlantMonitoring) TableName() string {
	return "plant_monitorings"
}
