package harvest

import "time"

type HarvestOutput struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	HarvestID uint    `gorm:"not null;index" json:"harvest_id"`
	Jenis     string  `gorm:"type:varchar(20);not null;index" json:"jenis"`
	Quantity  float64 `gorm:"type:numeric(12,2);not null" json:"quantity"`
	Status    string  `gorm:"type:varchar(50);not null;index" json:"status"`

	Harvest Harvest `gorm:"foreignKey:HarvestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (HarvestOutput) TableName() string {
	return "harvest_outputs"
}
