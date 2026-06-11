package notification

import "time"

type EmailMonitoring struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	RecipientEmail string     `gorm:"type:varchar(150);not null;index:idx_email_monitoring_recipient_template_date,priority:1" json:"recipient_email"`
	TemplateName   string     `gorm:"type:varchar(100);not null;index:idx_email_monitoring_recipient_template_date,priority:2" json:"template_name"`
	MonitorDate    time.Time  `gorm:"type:date;not null;index:idx_email_monitoring_recipient_template_date,priority:3" json:"monitor_date"`
	MessageCount   int        `gorm:"not null;default:0" json:"message_count"`
	FirstSentAt    *time.Time `json:"first_sent_at,omitempty"`
	LastSentAt     *time.Time `json:"last_sent_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (EmailMonitoring) TableName() string {
	return "email_monitorings"
}
