package harvest

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	maintenancemodels "pondok-tani-backend/models/maintenance"
	"time"
)

type HarvestRequest struct {
	ID uint `gorm:"primaryKey" json:"id"`

	PlantMonitoringID uint   `gorm:"not null;index" json:"plant_monitoring_id"`
	PlantBatchID      uint   `gorm:"not null;index" json:"plant_batch_id"`
	RequestedBy       string `gorm:"type:uuid;not null;index" json:"requested_by"`
	Status            string `gorm:"type:varchar(20);not null;default:pending;index" json:"status"`

	Phase       string `gorm:"type:varchar(50);not null" json:"phase"`
	Progress    int    `gorm:"not null" json:"progress"`
	HealthStatus string `gorm:"type:varchar(30);not null" json:"health_status"`
	Disease     string `gorm:"type:varchar(50)" json:"disease"`
	DiseaseNote string `gorm:"type:text" json:"disease_note"`
	AffectedCount int `json:"affected_count"`
	TotalPlants   int `json:"total_plants"`
	Note        string `gorm:"type:text" json:"note"`
	PhotoURL    string `gorm:"type:text" json:"photo_url"`
	MonitoringDate time.Time `json:"monitoring_date"`
	TotalQuantity float64 `gorm:"type:numeric(12,2);not null;default:0" json:"total_quantity"`
	WetQuantity   float64 `gorm:"type:numeric(12,2);not null;default:0" json:"wet_quantity"`
	DryQuantity   float64 `gorm:"type:numeric(12,2);not null;default:0" json:"dry_quantity"`

	ValidatedBy *string    `gorm:"type:uuid;index" json:"validated_by,omitempty"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
	HarvestID   *uint      `gorm:"index" json:"harvest_id,omitempty"`

	PlantMonitoring maintenancemodels.PlantMonitoring `gorm:"foreignKey:PlantMonitoringID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	PlantBatch      coremodels.PlantBatch             `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Requester       authmodels.User                   `gorm:"foreignKey:RequestedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Validator       authmodels.User                   `gorm:"foreignKey:ValidatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Harvest         Harvest                           `gorm:"foreignKey:HarvestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (HarvestRequest) TableName() string {
	return "harvest_requests"
}
