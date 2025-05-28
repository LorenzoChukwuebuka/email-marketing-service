package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/campaigns/dto"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"github.com/google/uuid"
)

type Service struct {
	store db.Store
}

func NewCampaignService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateCampaign(ctx context.Context, req *dto.CampaignDTO) (*dto.CampaignDTO, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserId,
		"template": req.TemplateId,
	})

	exists, err := s.store.CheckCampaignNameExists(ctx, db.CheckCampaignNameExistsParams{
		CompanyID: _uuid["company"],
		Name:      req.Name,
	})
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if !exists {
		return nil, common.ErrRecordExists
	}

	_, err = s.store.CreateCampaign(ctx, db.CreateCampaignParams{
		CompanyID:      _uuid["company"],
		Name:           req.Name,
		Subject:        sql.NullString{String: req.Subject, Valid: true},
		PreviewText:    sql.NullString{String: req.PreviewText, Valid: true},
		UserID:         _uuid["user"],
		SenderFromName: sql.NullString{String: req.SenderFromName, Valid: true},
		TemplateID:     uuid.NullUUID{UUID: _uuid["template"], Valid: true},
		RecipientInfo:  sql.NullString{String: req.RecipientInfo, Valid: true},
		Status:         sql.NullString{String: string(enums.CampaignStatus(enums.Draft)), Valid: true},
		TrackType:      sql.NullString{String: req.TrackType, Valid: true},
		Sender:         sql.NullString{String: req.Sender, Valid: true},
		ScheduledAt:    sql.NullTime{Time: req.ScheduledAt, Valid: true},
		HasCustomLogo:  sql.NullBool{Bool: req.HasCustomLogo, Valid: true},
	})

	return nil, nil
}
