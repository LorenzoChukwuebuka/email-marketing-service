package utils

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
)

// CreateNotification is a utility function to easily create a new notification
func CreateNotification(repo *repository.UserNotificationRepository, userId, title string) error {
	notification := &model.UserNotification{
		UserId:     userId,
		Title:      title,
		ReadStatus: false,
	}

	return repo.CreateNotification(notification)
}
