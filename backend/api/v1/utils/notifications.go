package utils

import (
	"email-marketing-service/api/v1/model"
	adminmodel "email-marketing-service/api/v1/model/admin"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
)

// CreateNotification is a utility function to easily create a new notification
func CreateNotification(repo *repository.UserNotificationRepository, userId string, title string) error {
	notification := &model.UserNotification{
		UserId:     userId,
		Title:      title,
		ReadStatus: false,
	}

	return repo.CreateNotification(notification)
}

func CreateAdminNotifications(repo *adminrepository.AdminNotificationRepository, userId, link, title string) error {
	notification := &adminmodel.AdminNotification{
		UserId:     userId,
		Link:       link,
		Title:      title,
		ReadStatus: false,
	}
	return repo.CreateNewNotification(notification)
}
