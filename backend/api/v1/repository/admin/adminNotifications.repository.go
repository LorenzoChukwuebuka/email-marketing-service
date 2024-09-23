package adminrepository

import (
	adminmodel "email-marketing-service/api/v1/model/admin"
	"time"

	"gorm.io/gorm"
)

type AdminNotificationRepository struct {
	DB *gorm.DB
}

func NewAdminNoficationRepository(db *gorm.DB) *AdminNotificationRepository {
	return &AdminNotificationRepository{
		DB: db,
	}
}

func (r *AdminNotificationRepository) CreateNewNotification(d *adminmodel.AdminNotification) error {
	if err := r.DB.Create(d).Error; err != nil {
		return err
	}
	return nil
}

func (r *AdminNotificationRepository) GetAlNotifications() ([]adminmodel.AdminNotificationResponse, error) {
	var notifications []adminmodel.AdminNotification
	if err := r.DB.Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	notificationResponses := make([]adminmodel.AdminNotificationResponse, 0)
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, adminmodel.AdminNotificationResponse{
			UUID:       notification.UUID,
			UserId:     notification.UserId,
			Title:      notification.Title,
			Link:       notification.Link,
			ReadStatus: notification.ReadStatus,
			CreatedAt:  notification.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  notification.UpdatedAt.Format(time.RFC3339),
			DeletedAt:  nil,
		})
	}

	return notificationResponses, nil
}
