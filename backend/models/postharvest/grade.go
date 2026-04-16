package postharvest

import "time"

type Grade struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	GradeName   string  `gorm:"type:varchar(40);not null;uniqueIndex" json:"grade_name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Grade) TableName() string {
	return "grades"
}
