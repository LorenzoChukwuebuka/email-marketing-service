package repository

import (
	"email-marketing-service/api/v1/model"
	"log"
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

	result := ls.DB.Create(&logEntry)
	if result.Error != nil {
		log.Printf("Error creating log entry: %v", result.Error)
		return result.Error
	}

	log.Printf("Log entry created successfully: UserID=%s, Action=%s, Details=%s", userID, action, details)
	return nil

}
