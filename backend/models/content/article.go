package content

import "time"

type Article struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(150);not null" json:"title"`
	Slug        string     `gorm:"type:varchar(180);not null;uniqueIndex" json:"slug"`
	Category    string     `gorm:"type:varchar(80);not null;default:'Umum'" json:"category"`
	Excerpt     string     `gorm:"type:text" json:"excerpt"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	ImageURL    string     `gorm:"type:text" json:"image_url"`
	Status      string     `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	Author      string     `gorm:"type:varchar(100);not null;default:'Admin'" json:"author"`
	PublishedAt *time.Time `gorm:"index" json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
