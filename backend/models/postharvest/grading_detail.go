package postharvest

import "time"

type GradingDetail struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	GradingBatchID uint    `gorm:"not null;index" json:"grading_batch_id"`
	GradeID        uint    `gorm:"not null;index" json:"grade_id"`
	Quantity       float64 `gorm:"type:numeric(12,2);not null" json:"quantity"`

	GradingBatch GradingBatch `gorm:"foreignKey:GradingBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Grade        Grade        `gorm:"foreignKey:GradeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (GradingDetail) TableName() string {
	return "grading_details"
}
