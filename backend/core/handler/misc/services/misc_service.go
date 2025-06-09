package services

import (
	"context"
	"database/sql"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"time"
)

type MiscService struct {
	store db.Store
}

func NewMiscService(store db.Store) *MiscService {
	return &MiscService{
		store: store,
	}
}

func (s *MiscService) TrackOpenCampaignEmails(ctx context.Context, campaignId string, email string, deviceType string, ipAddress string) error {
	htime := time.Now().UTC()

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
	})

	if err != nil {
		return common.ErrInvalidUUID
	}

	// Retrieve the existing email campaign result
	existingEmailResult, err := s.store.GetEmailCampaignResult(ctx, db.GetEmailCampaignResultParams{ID: _uuid["campaign"], RecipientEmail: email})
	if err != nil {
		return err
	}

	// Increment the OpenCount by 1
	var openCount int32 = 1

	if existingEmailResult != (db.EmailCampaignResult{}) {
		openCount = existingEmailResult.OpenCount.Int32 + 1
	}

	_, err = s.store.UpdateEmailCampaignResult(ctx, db.UpdateEmailCampaignResultParams{
		CampaignID:     _uuid["campaign"],
		RecipientEmail: email,
		OpenedAt:       sql.NullTime{Time: htime, Valid: true},
		Location:       sql.NullString{String: ipAddress, Valid: true},
		DeviceType:     sql.NullString{String: deviceType, Valid: true},
		OpenCount:      sql.NullInt32{Int32: openCount, Valid: true},
	})

	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}

func (s *MiscService) UnsubscribeFromCampaign(ctx context.Context, campaignId string, email string, companyId string) error {
	htime := time.Now().UTC()
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
		"company":  companyId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}
	err = s.store.ExecTx(ctx, func(q *db.Queries) error {

		_, err = q.UpdateEmailCampaignResult(ctx, db.UpdateEmailCampaignResultParams{
			CampaignID:     _uuid["campaign"],
			RecipientEmail: email,
			UnsubscribedAt: sql.NullTime{Time: htime, Valid: true},
		})

		if err != nil {
			return err
		}

		err = q.UpdateContactSubscription(ctx, db.UpdateContactSubscriptionParams{
			CompanyID:    _uuid["company"],
			Email:        email,
			IsSubscribed: sql.NullBool{Bool: false, Valid: true},
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}

func (s *MiscService) TrackClickedCampaignsEmails(ctx context.Context, campaignId string, email string) error {
	htime := time.Now().UTC()
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
	})

	if err != nil {
		return common.ErrInvalidUUID
	}

	// Retrieve the existing email campaign result
	existingEmailResult, err := s.store.GetEmailCampaignResult(ctx, db.GetEmailCampaignResultParams{ID: _uuid["campaign"], RecipientEmail: email})
	if err != nil {
		return err
	}

	var clickCount int32 = 1
	if existingEmailResult != (db.EmailCampaignResult{}) {
		clickCount = existingEmailResult.OpenCount.Int32 + 1
	}

	_, err = s.store.UpdateEmailCampaignResult(ctx, db.UpdateEmailCampaignResultParams{
		CampaignID:     _uuid["campaign"],
		RecipientEmail: email,
		ClickedAt:      sql.NullTime{Time: htime, Valid: true},
		ClickCount:     sql.NullInt32{Int32: clickCount, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *MiscService) GetPlans() (any, error) {
	return nil, nil
}
