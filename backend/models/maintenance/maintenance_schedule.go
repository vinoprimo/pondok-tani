package maintenance

import (
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type MaintenanceSchedule struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PlantBatchID  uint      `gorm:"not null;index" json:"plant_batch_id"`
	ActivityType  string    `gorm:"type:varchar(60);not null;index" json:"activity_type"`
	FrequencyDays uint      `gorm:"not null;default:0" json:"frequency_days"`
	NextDueDate   time.Time `gorm:"not null;index" json:"next_due_date"`
	Status        string    `gorm:"type:varchar(20);not null;default:pending;index" json:"status"`

	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MaintenanceSchedule) TableName() string {
	return "maintenance_schedules"
}
