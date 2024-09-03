package repository

import (
	"email-marketing-service/api/v1/model"
	"gorm.io/gorm"
)

type LogRepository struct {
	DB *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{DB: db}
}

func (ls *LogRepository) LogAction(userID, action, details string) error {
	logEntry := model.Log{
		UserID:  userID,
		Action:  action,
		Details: details,
	}

	return ls.DB.Create(&logEntry).Error
}
