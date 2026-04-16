package price

import (
	postharvestmodels "pondok-tani-backend/models/postharvest"
	"time"
)

type NationalPrice struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	GradeID       *uint     `gorm:"index" json:"grade_id,omitempty"`
	Province      *string   `gorm:"type:varchar(120)" json:"province,omitempty"`
	PricePerKg    float64   `gorm:"type:numeric(14,2);not null" json:"price_per_kg"`
	EffectiveDate time.Time `gorm:"not null;index" json:"effective_date"`
	Source        *string   `gorm:"type:varchar(150)" json:"source,omitempty"`

	Grade *postharvestmodels.Grade `gorm:"foreignKey:GradeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NationalPrice) TableName() string {
	return "national_prices"
}
