package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/auth/dto"
	"email-marketing-service/core/handler/user/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"time"
)

type UserService struct {
	store db.Store
}

func NewUserService(store db.Store) *UserService {
	return &UserService{
		store: store,
	}
}

const (
	DeletionGracePeriod   = 30 * 24 * time.Hour // 30 days
	StatusActive          = "active"
	StatusPendingDeletion = "pending_deletion"
	StatusDeleted         = "deleted"
)

func (s *UserService) GetUserNotifications(ctx context.Context, userID string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userID,
	})
	if err != nil {
		return nil, err
	}

	notifications, err := s.store.GetUserNotifications(ctx, _uuid["user"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}
	return notifications, nil
}

func (s *UserService) GetUserNotificationsLongPoll(ctx context.Context, userID string, sinceID *string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userID,
	})
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Timeout or cancellation - return empty array
			return []interface{}{}, nil

		case <-ticker.C:
			// Check for new notifications
			var notifications any
			var err error

			if sinceID != nil {
				// Parse the sinceID UUID
				sinceUUID, parseErr := common.ParseUUIDMap(map[string]string{
					"since": *sinceID,
				})
				if parseErr != nil {
					return nil, parseErr
				}

				// Get notifications created after the notification with sinceID
				notifications, err = s.store.GetUserNotificationsSinceID(ctx, db.GetUserNotificationsSinceIDParams{
					UserID: _uuid["user"],
					ID:     sinceUUID["since"],
				})
			} else {
				// First request - get all notifications
				notifications, err = s.store.GetUserNotifications(ctx, _uuid["user"])
			}

			if err != nil && err != sql.ErrNoRows {
				return nil, common.ErrFetchingRecord
			}

			// If we found new notifications, return immediately
			if notifications != nil && s.hasNotifications(notifications) {
				return notifications, nil
			}
		}
	}
}

// Helper to check if we have actual notifications
func (s *UserService) hasNotifications(data any) bool {
	switch v := data.(type) {
	case []interface{}:
		return len(v) > 0
	case []map[string]interface{}:
		return len(v) > 0
	default:
		return false
	}
}

func (s *UserService) GetUserDetails(ctx context.Context, userId string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userId,
	})
	if err != nil {
		return nil, err
	}

	userDetails, err := s.store.GetUserByID(ctx, _uuid["user"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	data := mapper.MapUserResponse(userDetails)
	return data, err
}

func (s *UserService) UpdateReadStatus(ctx context.Context, userId string) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}
	err = s.store.MarkAdminNotificationAsRead(ctx, _uuid["user"])
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) MarkUserForDeletion(ctx context.Context, userId string) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	_, err = s.store.MarkUserForDeletion(ctx, db.MarkUserForDeletionParams{
		ID:                  _uuid["user"],
		ScheduledDeletionAt: sql.NullTime{Time: time.Now().Add(DeletionGracePeriod), Valid: true},
		Status:              StatusPendingDeletion,
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) CancelUserDeletion(ctx context.Context, userId string) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	_, err = s.store.CancelUserDeletion(ctx, _uuid["user"])
	if err != nil {
		return nil
	}

	return nil
}

func (c *UserService) GetCurrentRunningSubscription(ctx context.Context, company_id string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": company_id,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	subscription, err := c.store.GetCurrentRunningSubscription(ctx, _uuid["company"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	emailUsage, err := c.store.GetCurrentEmailUsage(ctx, _uuid["company"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	return map[string]interface{}{
		"plan":           subscription.PlanName,
		"mailsPerDay":    emailUsage.EmailsLimit,
		"remainingMails": emailUsage.RemainingEmails.Int32,
	}, nil
}

func (s *UserService) EditUser(ctx context.Context, userId string, companyId string, req *dto.EditUserDTO) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":    userId,
		"company": companyId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.UpdateUserRecords(ctx, db.UpdateUserRecordsParams{
		ID:          _uuid["user"],
		Fullname:    sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		Phonenumber: sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
	})

	if err != nil {
		return err
	}

	//update company
	err = s.store.UpdateCompanyName(ctx, db.UpdateCompanyNameParams{
		ID:          _uuid["company"],
		Companyname: sql.NullString{String: req.Company, Valid: req.Company != ""},
	})

	if err != nil {
		return err
	}

	return nil
}
