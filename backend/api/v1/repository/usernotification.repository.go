package repository

import (
	"email-marketing-service/api/v1/model"
	//"errors"
	"time"

	"gorm.io/gorm"
)

type UserNotificationRepository struct {
	DB *gorm.DB
}

func NewUserNotificationRepository(db *gorm.DB) *UserNotificationRepository {
	return &UserNotificationRepository{
		DB: db,
	}
}

// CreateNotification creates a new notification for the user.
func (r *UserNotificationRepository) CreateNotification(d *model.UserNotification) error {
	if err := r.DB.Create(d).Error; err != nil {
		return err
	}
	return nil
}

// GetAllUserNotification retrieves all notifications for a specific user.
func (r *UserNotificationRepository) GetAllUserNotification(userId string) ([]model.UserNotificationResponse, error) {
    var notifications []model.UserNotification
    if err := r.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&notifications).Error; err != nil {
        return nil, err
    }

    notificationResponses := make([]model.UserNotificationResponse, 0)
    for _, notification := range notifications {
        notificationResponses = append(notificationResponses, model.UserNotificationResponse{
            UUID:       notification.UUID,
            UserId:     notification.UserId,
            Title:      notification.Title,
            ReadStatus: notification.ReadStatus,
            CreatedAt:  notification.CreatedAt.Format(time.RFC3339),
            UpdatedAt:  notification.UpdatedAt.Format(time.RFC3339),
            DeletedAt:  nil,
        })
    }

    return notificationResponses, nil
}

// UpdateReadStatus updates the read status of all notifications for a user to true.
func (r *UserNotificationRepository) UpdateReadStatus(userId string) error {
	result := r.DB.Model(&model.UserNotification{}).Where("user_id = ? AND read_status = ?", userId, false).Update("read_status", true)
	if result.Error != nil {
		return result.Error
	}

	// if result.RowsAffected == 0 {
	// 	return errors.New("no unread notifications found")
	// }

	return nil
}
