package adminservice

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"fmt"
)

type AdminCampaignService struct {
	AdminCampaignRepo *adminrepository.AdminCampaignRepository
}

func NewAdminCampaignService(adminCampaignRepo *adminrepository.AdminCampaignRepository) *AdminCampaignService {
	return &AdminCampaignService{
		AdminCampaignRepo: adminCampaignRepo,
	}
}

func (s *AdminCampaignService) GetAllUserCampaigns(userId, search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	campaigns, err := s.AdminCampaignRepo.GetAllUserCampaigns(userId, search, params)
	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to get user campaigns: %w", err)
	}
	return campaigns, nil
}

func (s *AdminCampaignService) GetAUserCampaign(uuid string) (*model.CampaignResponse, error) {
	campaign, err := s.AdminCampaignRepo.GetAUserCampaign(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user campaign: %w", err)
	}
	return campaign, nil
}

func (s *AdminCampaignService) DeleteUserCampaign(uuid string) error {
	err := s.AdminCampaignRepo.DeleteUserCampaign(uuid)
	if err != nil {
		return fmt.Errorf("failed to delete user campaign: %w", err)
	}
	return nil
}

func (s *AdminCampaignService) SuspendUserCampaign(uuid string) error {
	err := s.AdminCampaignRepo.SuspendUserCampaign(uuid)
	if err != nil {
		return fmt.Errorf("failed to suspend user campaign: %w", err)
	}
	return nil
}
