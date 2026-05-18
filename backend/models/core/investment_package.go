package core

import "time"

type InvestmentPackage struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	PackageName string   `gorm:"type:varchar(120);not null;index" json:"package_name"`
	Description *string  `gorm:"type:text" json:"description,omitempty"`
	MinQuantity uint     `gorm:"not null;default:1" json:"min_quantity"`
	Price       float64  `gorm:"type:numeric(14,2);not null" json:"price"`
	Status      string   `gorm:"type:varchar(20);not null;default:active;index" json:"status"`
	Roi         string   `gorm:"type:varchar(50);default:''" json:"roi"`
	Duration    string   `gorm:"type:varchar(50);default:''" json:"duration"`
	Benefits    []string `gorm:"type:jsonb;serializer:json" json:"benefits"`
	IsPopular   bool     `gorm:"default:false" json:"is_popular"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (InvestmentPackage) TableName() string {
	return "investment_packages"
}
