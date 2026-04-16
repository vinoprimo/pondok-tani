package core

import (
	authmodels "pondok-tani-backend/models/auth"
	"time"
)

type Investment struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           string     `gorm:"type:uuid;not null;index" json:"user_id"`
	PackageID        uint       `gorm:"not null;index" json:"package_id"`
	Amount           float64    `gorm:"type:numeric(14,2);not null" json:"amount"`
	Status           string     `gorm:"type:varchar(20);not null;default:pending;index" json:"status"`
	InvestmentDate   time.Time  `gorm:"not null" json:"investment_date"`
	ExpectedReturnAt *time.Time `json:"expected_return_at,omitempty"`

	User    authmodels.User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Package InvestmentPackage `gorm:"foreignKey:PackageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Investment) TableName() string {
	return "investments"
}
