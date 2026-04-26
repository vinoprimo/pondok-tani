package maintenance

import (
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type MaintenanceSchedule struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        string    `gorm:"type:uuid;index" json:"user_id"`
	PlantBatchID  uint      `gorm:"not null;index" json:"plant_batch_id"`
	ActivityType  string    `gorm:"type:varchar(60);not null;index" json:"activity_type"`
	FrequencyDays uint      `gorm:"not null;default:0" json:"frequency_days"`
	NextDueDate   time.Time `gorm:"not null;index" json:"next_due_date"`
	Status        string    `gorm:"type:varchar(20);not null;default:pending;index" json:"status"`

	User       authmodels.User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	PlantBatch coremodels.PlantBatch `gorm:"foreignKey:PlantBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Activities []MaintenanceActivity `gorm:"foreignKey:ScheduleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MaintenanceSchedule) TableName() string {
	return "maintenance_schedules"
}
