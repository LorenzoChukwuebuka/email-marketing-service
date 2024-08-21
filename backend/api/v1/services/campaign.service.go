package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
)

type CampaignService struct {
	CampaignRepo *repository.CampaignRepository
	ContactRepo  *repository.ContactRepository
}

func NewCampaignService(campaignRepo *repository.CampaignRepository, contactRepo *repository.ContactRepository) *CampaignService {
	return &CampaignService{
		CampaignRepo: campaignRepo,
		ContactRepo:  contactRepo,
	}
}

func (s *CampaignService) CreateCampaign(d *dto.CampaignDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid  data: %w", err)
	}

	campaignModel := &model.Campaign{Name: d.Name, UserId: d.UserId, Status: model.CampaignStatus(dto.Draft), SenderFromName: d.SenderFromName}

	campaignExist, err := s.CampaignRepo.CampaignExists(campaignModel)

	if err != nil {
		return nil, err
	}

	if campaignExist {
		return nil, fmt.Errorf("campaign already exists")
	}

	saveCampaign, err := s.CampaignRepo.CreateCampaign(campaignModel)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"campaignId": saveCampaign,
	}, nil
}

func (s *CampaignService) GetAllCampaigns(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	campaignRepo, err := s.CampaignRepo.GetAllCampaigns(userId, paginationParams)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	if campaignRepo.TotalCount == 0 {
		return repository.PaginatedResult{}, nil
	}
	return campaignRepo, nil
}

func (s *CampaignService) GetScheduledCampaigns(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	campaignRepo, err := s.CampaignRepo.GetAllCampaigns(userId, paginationParams)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	if campaignRepo.TotalCount == 0 {
		return repository.PaginatedResult{}, nil
	}
	return campaignRepo, nil
}

func (s *CampaignService) GetSingleCampaign(userId string, campaignId string) (*model.CampaignResponse, error) {

	campaignRepo, err := s.CampaignRepo.GetSingleCampaign(userId, campaignId)

	if err != nil {
		return nil, err
	}

	if campaignRepo == nil {
		return nil, nil
	}

	return campaignRepo, nil
}

func (s *CampaignService) UpdateCampaign(d *dto.CampaignDTO) error {

	campaignModel := &model.Campaign{
		UUID:           d.UUID,
		Name:           d.Name,
		Subject:        d.Subject,
		PreviewText:    d.PreviewText,
		UserId:         d.UserId,
		SenderFromName: d.SenderFromName,
		TemplateId:     d.TemplateId,
		SentTemplateId: d.SentTemplateId,
		RecipientInfo:  d.RecipientInfo,
		IsPublished:    d.IsPublished,
		Status:         model.CampaignStatus(d.Status),
		TrackType:      model.Track,
		IsArchived:     d.IsArchived,
		SentAt:         d.SentAt,
		Template:       d.Template,
		ScheduledAt:    d.ScheduledAt,
		HasCustomLogo:  d.HasCustomLogo,
	}

	if err := s.CampaignRepo.UpdateCampaign(campaignModel); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) AddOrEditCampaignGroup(d *dto.CampaignGroupDTO) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid  data: %w", err)
	}

	getCampaign, err := s.CampaignRepo.GetSingleCampaign(d.UserId, d.CampaignId)
	if err != nil {
		return err
	}

	getContactGroup, err := s.ContactRepo.GetASingleGroup(d.UserId, d.GroupId)
	if err != nil {
		return err
	}

	cgpModel := &model.CampaignGroup{CampaignId: getCampaign.ID, GroupId: getContactGroup.ID}

	if err := s.CampaignRepo.AddOrEditCampaignGroup(cgpModel); err != nil {
		return err
	}

	return nil
}
