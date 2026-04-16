package postharvest

import "time"

type GradingBatch struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DryingProcessID uint      `gorm:"not null;index" json:"drying_process_id"`
	BatchCode       string    `gorm:"type:varchar(60);not null;uniqueIndex" json:"batch_code"`
	GradingDate     time.Time `gorm:"not null;index" json:"grading_date"`
	Status          string    `gorm:"type:varchar(30);not null;default:completed;index" json:"status"`

	DryingProcess DryingProcess `gorm:"foreignKey:DryingProcessID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (GradingBatch) TableName() string {
	return "grading_batches"
}
