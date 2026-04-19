package core

import (
	"time"
)

type Investment struct {
	ID             uint              `gorm:"primaryKey" json:"id"`
	UserID         string            `gorm:"type:uuid;not null;index:idx_investments_user_package,unique" json:"user_id"`
	PackageID      uint              `gorm:"not null;index:idx_investments_user_package,unique" json:"package_id"`
	Amount         float64           `gorm:"type:numeric(14,2);not null" json:"amount"`
	Status         string            `gorm:"type:varchar(20);not null;default:pending;index" json:"status"`
	InvestmentDate time.Time         `gorm:"not null" json:"investment_date"`
	Package        InvestmentPackage `gorm:"foreignKey:PackageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"package,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Investment) TableName() string {
	return "investments"
}
