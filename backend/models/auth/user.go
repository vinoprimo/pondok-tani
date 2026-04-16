package auth

import "time"

type User struct {
	ID          string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Email       string  `gorm:"type:varchar(150);not null;uniqueIndex" json:"email"`
	Password    string  `gorm:"type:text;not null" json:"-"`
	Role        string  `gorm:"type:varchar(20);not null;default:investor;index" json:"role"`
	PhoneNumber *string `gorm:"type:varchar(30)" json:"phone_number,omitempty"`
	Address     *string `gorm:"type:text" json:"address,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
