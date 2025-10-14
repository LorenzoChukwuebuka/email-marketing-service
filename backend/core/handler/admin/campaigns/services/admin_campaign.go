package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/admin/campaigns/dto"
	campaignDTO "email-marketing-service/core/handler/campaigns/dto"
	campaignMapper "email-marketing-service/core/handler/campaigns/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"fmt"
	"github.com/google/uuid"
)

type AdminCampaignService struct {
	store db.Store
}

func NewAdminCampaignService(store db.Store) *AdminCampaignService {
	return &AdminCampaignService{store: store}
}

type CampaignWithGroupsResponse struct {
	*campaignDTO.CampaignResponseDTO
	Groups []*campaignDTO.GetCampaignContactGroupsResponse `json:"groups"`
}

func (s *AdminCampaignService) GetAllUserCampaigns(ctx context.Context, req *dto.AdminFetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":    req.UserID,
		"company": req.CompanyID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaigns, err := s.store.ListCampaignsByCompanyID(ctx, db.ListCampaignsByCompanyIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
		Column5:   req.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching record : %w", err)
	}

	count_campaigns, err := s.store.GetCampaignCounts(ctx, db.GetCampaignCountsParams{
		UserID:    _uuid["user"],
		CompanyID: _uuid["company"],
	})
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	response := campaignMapper.MapCampaignResponses(campaigns)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}
	data := common.Paginate(int(count_campaigns), items, req.Offset, req.Limit)
	return data, nil
}

func (s *AdminCampaignService) GetSingleCampaign(ctx context.Context, req *dto.AdminFetchCampaignDTO) (any, error) {
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

	//get contact groups
	campaign_group, err := s.store.GetCampaignContactGroups(ctx, _uuid["campaign"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	// Map the campaign data to the response DTO
	groupData := campaignMapper.MapCampaignGroups(campaign_group)
	campaignData := campaignMapper.MapGetCampaignResponse(campaign)

	// Fetch template separately if template_id exists
	if campaign.TemplateID.Valid && campaign.TemplateID.UUID != uuid.Nil {
		template, err := s.store.GetTemplateByIDWithoutType(ctx, db.GetTemplateByIDWithoutTypeParams{
			ID:     campaign.TemplateID.UUID,
			UserID: uuid.NullUUID{UUID: _uuid["user"]},
		})

		if err != nil {
			if err == sql.ErrNoRows {
				// Template not found, set to nil
				campaignData.Template = nil
			} else {
				return nil, fmt.Errorf("error fetching template: %w", err)
			}
		} else {
			// Map the template data
			templateData := campaignMapper.MapTemplateFromSeparateQuery(template)
			campaignData.Template = templateData
		}
	}

	return &CampaignWithGroupsResponse{
		CampaignResponseDTO: campaignData,
		Groups:              groupData,
	}, nil
}

func (s *AdminCampaignService) GetAllRecipientsForACampaign(ctx context.Context, campaignId string, companyId string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
		"company":  companyId,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	result, err := s.store.GetEmailCampaignResultsByCampaign(ctx, db.GetEmailCampaignResultsByCampaignParams{
		CampaignID: _uuid["campaign"],
		CompanyID:  _uuid["company"],
	})
	if err != nil {
		return nil, err
	}

	data := campaignMapper.MapCampaignEmailResponse(result)
	return data, nil
}

func (s *AdminCampaignService) SuspendCampaign(ctx context.Context, req *dto.AdminFetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":     req.UserID,
		"campaign": req.CampaignID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	err = s.store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
		ID:     _uuid["campaign"],
		UserID: _uuid["user"],
		Status: sql.NullString{String: string(enums.Suspended), Valid: true},
	})

	if err != nil {
		return nil, common.ErrSuspendingCampaign
	}

	return nil, nil
}

func (s *AdminCampaignService) UnsuspendCampaign(ctx context.Context, req *dto.AdminFetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserID,
		"campaign": req.CampaignID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	err = s.store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
		ID:     _uuid["campaign"],
		UserID: _uuid["user"],
		Status: sql.NullString{String: string(enums.Draft), Valid: true},
	})
	if err != nil {
		return nil, common.ErrUnsuspendingCampaign
	}

	return nil, nil
}
