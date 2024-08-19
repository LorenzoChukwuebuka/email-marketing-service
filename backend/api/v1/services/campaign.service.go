package services

import "email-marketing-service/api/v1/dto"

type CampaignService struct{}

func NewCampaignService() *CampaignService {
	return &CampaignService{}
}

func (s *CampaignService) CreateCampaign(d *dto.CampaignDTO) (map[string]interface{},error) {
	return nil,nil
}
