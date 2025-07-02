package services

import (
	"context"
	"database/sql"
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

// func (s *UserService) EditUser(ctx context.Context, userId string) error {
// 	_uuid, err := common.ParseUUIDMap(map[string]string{
// 		"user": userId,
// 	})
// 	if err != nil {
// 		return common.ErrInvalidUUID
// 	}

// 	_, err = s.store.EditUser(ctx, _uuid["user"])
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
