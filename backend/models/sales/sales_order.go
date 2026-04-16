package sales

import (
	authmodels "pondok-tani-backend/models/auth"
	"time"
)

type SalesOrder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderNumber string    `gorm:"type:varchar(60);not null;uniqueIndex" json:"order_number"`
	SalesDate   time.Time `gorm:"not null;index" json:"sales_date"`
	Buyer       string    `gorm:"type:varchar(150);not null" json:"buyer"`
	TotalAmount float64   `gorm:"type:numeric(14,2);not null;default:0" json:"total_amount"`
	Status      string    `gorm:"type:varchar(20);not null;default:draft;index" json:"status"`
	CreatedBy   string    `gorm:"type:uuid;not null;index" json:"created_by"`

	Creator authmodels.User `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SalesOrder) TableName() string {
	return "sales_orders"
}
