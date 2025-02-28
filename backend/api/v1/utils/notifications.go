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

// // Functional Option for UserNotification
// type UserNotificationOption func(*model.UserNotification)

// // Option function to set the read status
// func WithReadStatus(readStatus bool) UserNotificationOption {
// 	return func(n *model.UserNotification) {
// 		n.ReadStatus = readStatus
// 	}
// }

// // Option function to set an additional field
// func WithAdditionalField(field string) UserNotificationOption {
// 	return func(n *model.UserNotification) {
// 		n.AdditionalField = field // Assume UserNotification has this field
// 	}
// }

// // CreateNotification creates a UserNotification with functional options
// func CreateNotification(repo *repository.UserNotificationRepository, userId, title string, opts ...UserNotificationOption) error {
// 	notification := &model.UserNotification{
// 		UserId:     userId,
// 		Title:      title,
// 		ReadStatus: false, // Default value
// 	}

// 	// Apply all options
// 	for _, opt := range opts {
// 		opt(notification)
// 	}

// 	return repo.CreateNotification(notification)
// }

// // Functional Option for AdminNotification
// type AdminNotificationOption func(*adminmodel.AdminNotification)

// // Option function to set the admin read status
// func WithAdminReadStatus(readStatus bool) AdminNotificationOption {
// 	return func(n *adminmodel.AdminNotification) {
// 		n.ReadStatus = readStatus
// 	}
// }

// // Option function to set the admin link
// func WithAdminLink(link string) AdminNotificationOption {
// 	return func(n *adminmodel.AdminNotification) {
// 		n.Link = link
// 	}
// }

// // CreateAdminNotifications creates an AdminNotification with functional options
// func CreateAdminNotifications(repo *adminrepository.AdminNotificationRepository, userId, title string, opts ...AdminNotificationOption) error {
// 	notification := &adminmodel.AdminNotification{
// 		UserId:     userId,
// 		Title:      title,
// 		ReadStatus: false, // Default value
// 	}

// 	// Apply all options
// 	for _, opt := range opts {
// 		opt(notification)
// 	}

// 	return repo.CreateNewNotification(notification)
// }
