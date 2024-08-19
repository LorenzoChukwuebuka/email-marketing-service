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
}

func NewCampaignService(campaignRepo *repository.CampaignRepository) *CampaignService {
	return &CampaignService{
		CampaignRepo: campaignRepo,
	}
}

func (s *CampaignService) CreateCampaign(d *dto.CampaignDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid  data: %w", err)
	}

	campaignModel := &model.Campaign{Name: d.Name, UserId: d.UserId, Status: model.CampaignStatus(dto.Draft)}

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

func (s *CampaignService) GetScheduledCampaigns() error { return nil }

func (s *CampaignService) GetDraftCampaigns() error { return nil }

func (s *CampaignService) GetSentCampaigns() error { return nil }
