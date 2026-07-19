package maintenance

import (
	"errors"
	"pondok-tani-backend/config"
	"pondok-tani-backend/models/maintenance"
	"strings"
)

func GetAllActivityTypes(page, limit int, search string) ([]maintenance.ActivityType, int64, error) {
	var activityTypes []maintenance.ActivityType
	var total int64
	query := config.DB.Model(&maintenance.ActivityType{})

	if search != "" {
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Order("name ASC").Find(&activityTypes).Error; err != nil {
		return nil, 0, err
	}
	return activityTypes, total, nil
}

func GetActivityTypeByID(id uint) (*maintenance.ActivityType, error) {
	var activityType maintenance.ActivityType
	if err := config.DB.First(&activityType, id).Error; err != nil {
		return nil, errors.New("activity type not found")
	}
	return &activityType, nil
}

func CreateActivityType(activityType *maintenance.ActivityType) error {
	return config.DB.Create(activityType).Error
}

func UpdateActivityType(id uint, data map[string]interface{}) error {
	return config.DB.Model(&maintenance.ActivityType{}).Where("id = ?", id).Updates(data).Error
}

func DeleteActivityType(id uint) error {
	return config.DB.Delete(&maintenance.ActivityType{}, id).Error
}
