package maintenance

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type MaintenanceActivity struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ScheduleID       uint      `gorm:"not null;index" json:"schedule_id"`
	PlantBatchID     uint      `gorm:"not null;index" json:"plant_batch_id"`
	ActivityType     string    `gorm:"type:varchar(60);not null;index" json:"activity_type"`
	Description      *string   `gorm:"type:text" json:"description,omitempty"`
	ActivityDate     time.Time `gorm:"not null;index" json:"activity_date"`
	ValidationStatus string    `gorm:"type:varchar(20);not null;default:pending;index" json:"validation_status"`
	ValidationNotes  *string   `gorm:"type:text" json:"validation_notes,omitempty"`
	ValidatedBy      *string   `gorm:"type:uuid;index" json:"validated_by,omitempty"`

	Schedule   MaintenanceSchedule   `gorm:"foreignKey:ScheduleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Validator  *authmodels.User      `gorm:"foreignKey:ValidatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MaintenanceActivity) TableName() string {
	return "maintenance_activities"
}
