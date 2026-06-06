package notification

import "time"

type EmailLog struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         *string    `gorm:"type:uuid;index" json:"user_id,omitempty"`
	TemplateID     *uint      `gorm:"index" json:"template_id,omitempty"`
	RecipientEmail string     `gorm:"type:varchar(150);not null;index" json:"recipient_email"`
	Subject        string     `gorm:"type:varchar(255);not null" json:"subject"`
	Body           string     `gorm:"type:text;not null" json:"body"`
	Status         string     `gorm:"type:varchar(30);not null;default:'pending';index" json:"status"`
	ErrorMessage   *string    `gorm:"type:text" json:"error_message,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Template EmailTemplate `gorm:"foreignKey:TemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"template,omitempty"`
}

func (EmailLog) TableName() string {
	return "email_logs"
}
