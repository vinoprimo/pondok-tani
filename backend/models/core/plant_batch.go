package core

import "time"

type PlantBatch struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	InvestmentID uint      `gorm:"not null;index" json:"investment_id"`
	BatchCode    string    `gorm:"type:varchar(60);not null;uniqueIndex" json:"batch_code"`
	PlantingDate time.Time `gorm:"not null;index" json:"planting_date"`
	Location     string    `gorm:"type:varchar(150);not null;default:''" json:"location"`
	SeedCount    uint      `gorm:"not null;default:0" json:"seed_count"`
	LandArea     float64   `gorm:"type:numeric(12,2);not null;default:0" json:"land_area"`
	PhotoURL     string    `gorm:"type:text" json:"photo_url"`
	Status       string    `gorm:"type:varchar(30);not null;default:active;index" json:"status"`

	Investment Investment `gorm:"foreignKey:InvestmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PlantBatch) TableName() string {
	return "plant_batches"
}
