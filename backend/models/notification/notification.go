package notification

import (
	authmodels "pondok-tani-backend/models/auth"
	"time"
)

type Notification struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   string `gorm:"type:uuid;not null;index" json:"user_id"`
	UserName string `gorm:"type:varchar(100)" json:"user_name"`
	Type     string `gorm:"type:varchar(40);not null;index" json:"type"`
	Subject  string `gorm:"type:varchar(255)" json:"subject"`
	Message  string `gorm:"type:text;not null" json:"message"`
	IsRead   bool   `gorm:"not null;default:false;index" json:"is_read"`

	User authmodels.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
