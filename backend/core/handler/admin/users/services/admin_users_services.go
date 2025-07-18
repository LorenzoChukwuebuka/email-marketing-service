package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/admin/users/dto"
	"email-marketing-service/core/handler/admin/users/mapper"

	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"email-marketing-service/internal/mailer"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type AdminUsersServices struct {
	store db.Store
}

func NewAdminUsersServices(store db.Store) *AdminUsersServices {
	return &AdminUsersServices{
		store: store,
	}
}

const (
	DeletionGracePeriod   = 30 * 24 * time.Hour // 30 days
	StatusActive          = "active"
	StatusPendingDeletion = "pending_deletion"
	StatusDeleted         = "deleted"
)

var (
	otpLength = 8
)

func (s *AdminUsersServices) GetAllUsers(ctx context.Context, req *dto.AdminFetchUserDTO) (any, error) {
	result, err := s.store.GetAllUsers(ctx, db.GetAllUsersParams{
		Column1: req.Search,
		Offset:  int32(req.Offset),
		Limit:   int32(req.Limit),
	})

	if err != nil {
		return nil, common.ErrFetchingUser
	}

	total_count, err := s.store.CountAllUsers(ctx)
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	adapters := make([]mapper.GetAllUsersAdapter, len(result))
	for i, u := range result {
		adapters[i] = mapper.GetAllUsersAdapter(u)
	}

	response := mapper.MapAdminUsers(adapters)
	items := make([]any, len(response))

	for i, v := range response {
		items[i] = v
	}

	data := common.Paginate(int(total_count), items, req.Offset, req.Limit)
	return data, nil
}

func (s *AdminUsersServices) GetVerifiedUsers(ctx context.Context, req *dto.AdminFetchUserDTO) (any, error) {
	result, err := s.store.GetVerifiedUsers(ctx, db.GetVerifiedUsersParams{
		Column1: req.Search,
		Offset:  int32(req.Offset),
		Limit:   int32(req.Limit),
	})

	if err != nil {
		return nil, common.ErrFetchingUser
	}

	total_count, err := s.store.CountVerifiedUsers(ctx)
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	adapters := make([]mapper.GetVerifiedUsersAdapter, len(result))
	for i, u := range result {
		adapters[i] = mapper.GetVerifiedUsersAdapter(u)
	}

	mapped := mapper.MapAdminUsers(adapters)
	items := make([]any, len(mapped))

	for i, v := range mapped {
		items[i] = v
	}

	data := common.Paginate(int(total_count), items, req.Offset, req.Limit)
	return data, nil
}

func (s *AdminUsersServices) GetUnverifiedUsers(ctx context.Context, req *dto.AdminFetchUserDTO) (any, error) {
	result, err := s.store.GetUnVerifiedUsers(ctx, db.GetUnVerifiedUsersParams{
		Column1: req.Search,
		Offset:  int32(req.Offset),
		Limit:   int32(req.Limit),
	})

	if err != nil {
		return nil, common.ErrFetchingUser
	}

	total_count, err := s.store.CountUnVerifiedUsers(ctx)
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	adapters := make([]mapper.GetUnVerifiedUsersAdapter, len(result))
	for i, u := range result {
		adapters[i] = mapper.GetUnVerifiedUsersAdapter(u)
	}

	mapped := mapper.MapAdminUsers(adapters)
	items := make([]any, len(mapped))

	for i, v := range mapped {
		items[i] = v
	}

	data := common.Paginate(int(total_count), items, req.Offset, req.Limit)
	return data, nil
}

func (s *AdminUsersServices) VerifyUser(ctx context.Context, userId string) error {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}
	err = s.store.VerifyUser(ctx, uuiduserId)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminUsersServices) BlockUser(ctx context.Context, userId string) error {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.BlockUser(ctx, uuiduserId)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminUsersServices) UnblockUser(ctx context.Context, userId string) error {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}
	_, err = s.store.UnblockUser(ctx, uuiduserId)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminUsersServices) DeleteUser(ctx context.Context, userId string) error {
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

func (s *AdminUsersServices) GetUserByID(ctx context.Context, userId string) (any, error) {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	user, err := s.store.GetUserByID(ctx, uuiduserId)
	if err != nil {
		return nil, common.ErrFetchingUser
	}

	data := mapper.MapUserResponse(user)
	return data, err
}

func (s *AdminUsersServices) GetUserStats(ctx context.Context, userId string) (any, error) {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	var stats dto.AdminUserStats

	stats.TotalContacts, err = s.store.CountUserContacts(ctx, uuiduserId)
	if err != nil {
		return stats, fmt.Errorf("failed to count contacts: %w", err)
	}

	stats.TotalCampaigns, err = s.store.CountUserCampaigns(ctx, uuiduserId)
	if err != nil {
		return stats, fmt.Errorf("failed to count campaigns: %w", err)
	}

	stats.TotalTemplates, err = s.store.CountUserTemplates(ctx, uuiduserId)
	if err != nil {
		return stats, fmt.Errorf("failed to count templates: %w", err)
	}

	stats.TotalCampaignsSent, err = s.store.CountUserCampaignsSent(ctx, uuiduserId)
	if err != nil {
		return stats, fmt.Errorf("failed to count campaigns sent: %w", err)
	}

	stats.TotalGroups, err = s.store.CountUserGroups(ctx, uuiduserId)
	if err != nil {
		return stats, fmt.Errorf("failed to count groups: %w", err)
	}

	return stats, nil
}

func (s *AdminUsersServices) SendEmailToUsers(ctx context.Context, req *dto.AdminEmailLogDTO) error {
	if err := helper.ValidateData(req); err != nil {
		return errors.Join(common.ErrValidatingRequest, err)
	}

	users, err := s.store.GetAllVerifiedUserEmails(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	errChan := make(chan error, len(users))

	for _, user := range users {
		wg.Add(1)

		go func(email string) {
			defer wg.Done()

			print(user.Email)

			helper.AsyncSendMail(req.Subject, user.Email, req.Content, "info@crabmailer.com", nil, &wg)

		}(user.Email)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect any errors from the channel
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	return nil
}

func (s *AdminUsersServices) ResendUserVerificationMail(ctx context.Context, userId string) error {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	user, err := s.store.GetUserByID(ctx, uuiduserId)
	if err != nil {
		return err
	}
	if user.Status == StatusActive {
		return common.ErrUserAlreadyVerified
	}
	if user.Status == StatusPendingDeletion {
		return common.ErrUserPendingDeletion
	}
	if user.Status == StatusDeleted {
		return common.ErrUserDeleted
	}

	if user.Verified {
		return common.ErrUserAlreadyVerified
	}

	token := helper.GenerateOTP(otpLength)

	_, err = s.store.CreateOTP(ctx, db.CreateOTPParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error creating OTP: %w", err)
	}

	go mailer.NewEmailService().SignUpMail(user.Email, user.Fullname, user.ID, token)

	return nil
}


func (s *AdminUsersServices) GetAllUserContacts(ctx context.Context, userId string) (any, error) {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	fmt.Print(uuiduserId)
	return nil, nil
}

func (s *AdminUsersServices) GetAllUserTemplates(ctx context.Context, userId string) (any, error) {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	fmt.Print(uuiduserId)
	return nil, nil
}

func (s *AdminUsersServices) GetAllUserGroups(ctx context.Context, userId string) (any, error) {
	uuiduserId, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	fmt.Print(uuiduserId)
	return nil, nil
}
