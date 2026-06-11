package notification

import "time"

type EmailTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Subject   string    `gorm:"type:varchar(255);not null" json:"subject"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	IsActive  bool      `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EmailTemplate) TableName() string {
	return "email_templates"
}
