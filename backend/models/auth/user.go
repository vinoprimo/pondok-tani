package auth

import (
	coremodels "pondok-tani-backend/models/core"
	"time"
)

type User struct {
	ID                 string                        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name               string                        `gorm:"type:varchar(100);not null" json:"name"`
	Email              string                        `gorm:"type:varchar(150);not null;uniqueIndex" json:"email"`
	Password           string                        `gorm:"type:text;not null" json:"-"`
	Role               string                        `gorm:"type:varchar(20);not null;default:investor;index" json:"role"`
	PhoneNumber        *string                       `gorm:"type:varchar(30)" json:"phone_number,omitempty"`
	Address            *string                       `gorm:"type:text" json:"address,omitempty"`
	FCMToken           *string                       `gorm:"type:text" json:"fcm_token,omitempty"`
	SelectedPackageID  *uint                         `gorm:"index" json:"selected_package_id,omitempty"`
	SelectedPackage    *coremodels.InvestmentPackage `gorm:"foreignKey:SelectedPackageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"selected_package,omitempty"`
	Investments        []coremodels.Investment       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"investments,omitempty"`
	PackageStatus      string                        `gorm:"type:varchar(20);not null;default:none;index" json:"package_status"`
	PackageSelectedAt  *time.Time                    `json:"package_selected_at,omitempty"`
	PackageActivatedAt *time.Time                    `json:"package_activated_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
