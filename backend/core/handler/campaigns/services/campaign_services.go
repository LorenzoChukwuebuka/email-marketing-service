package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/campaigns/dto"
	"email-marketing-service/core/handler/campaigns/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"fmt"
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
		"company": req.CompanyID,
		"user":    req.UserId,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	exists, err := s.store.CheckCampaignNameExists(ctx, db.CheckCampaignNameExistsParams{
		CompanyID: _uuid["company"],
		Name:      req.Name,
	})
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if exists {
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

	if err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}

	return req, nil
}

func (s *Service) GetAllCampaigns(ctx context.Context, req *dto.FetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
		"user":    req.UserID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaigns, err := s.store.ListCampaignsByCompanyID(ctx, db.ListCampaignsByCompanyIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
		Column5:   req.SearchQuery,
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	count_campaigns, err := s.store.GetCampaignCounts(ctx, db.GetCampaignCountsParams{
		UserID:    _uuid["user"],
		CompanyID: _uuid["company"],
	})

	if err != nil {
		return nil, common.ErrFetchingCount
	}

	response := mapper.MapCampaignResponses(campaigns)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}
	data := common.Paginate(int(count_campaigns), items, req.Offset, req.Limit)
	return data, nil
}

func (s *Service) GetSingleCampaign(ctx context.Context, req *dto.FetchCampaignDTO) (*dto.CampaignResponseDTO, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserID,
		"campaign": req.CampaignID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaign, err := s.store.GetCampaignByID(ctx, db.GetCampaignByIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		ID:        _uuid["campaign"],
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	data := mapper.MapCampaignResponse(db.ListCampaignsByCompanyIDRow(campaign))
	return data, nil
}

func (s *Service) UpdateCampaign(ctx context.Context, req *dto.CampaignDTO, campaignId string) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserId,
		"campaign": campaignId,
		"template": req.TemplateId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	_, err = s.store.UpdateCampaign(ctx, db.UpdateCampaignParams{
		ID:             _uuid["campaign"],
		UserID:         _uuid["user"],
		Name:           req.Name,
		Subject:        sql.NullString{String: req.Subject, Valid: true},
		PreviewText:    sql.NullString{String: req.PreviewText, Valid: true},
		SenderFromName: sql.NullString{String: req.SenderFromName, Valid: true},
		TemplateID:     uuid.NullUUID{UUID: _uuid["template"], Valid: true},
		RecipientInfo:  sql.NullString{String: req.RecipientInfo, Valid: true},
		Status:         sql.NullString{String: string(req.Status), Valid: true},
		TrackType:      sql.NullString{String: req.TrackType, Valid: true},
		Sender:         sql.NullString{String: req.Sender, Valid: true},
		ScheduledAt:    sql.NullTime{Time: req.ScheduledAt, Valid: true},
		HasCustomLogo:  sql.NullBool{Bool: req.HasCustomLogo, Valid: true},
		IsPublished:    sql.NullBool{Bool: req.IsPublished, Valid: true},
		IsArchived:     sql.NullBool{Bool: req.IsArchived, Valid: true},
	})

	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}
