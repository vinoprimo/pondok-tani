package price

import (
	postharvestmodels "pondok-tani-backend/models/postharvest"
	"time"
)

type NationalPrice struct {
	NationalPriceID uint      `gorm:"primaryKey;column:national_price_id" json:"national_price_id"`
	GradeID         *uint     `gorm:"column:grade_id;index" json:"grade_id,omitempty"`
	Province        string    `gorm:"type:varchar(120);not null" json:"province"`
	PricePerKg      float64   `gorm:"type:numeric(14,2);not null" json:"price_per_kg"`
	EffectiveDate   time.Time `gorm:"not null;index" json:"effective_date"`
	HarvestType     string    `gorm:"type:varchar(80);not null;index" json:"harvest_type"`

	Grade *postharvestmodels.Grade `gorm:"foreignKey:GradeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"grade,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NationalPrice) TableName() string {
	return "national_prices"
}
