package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type CampaignRepository struct {
	DB *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{
		DB: db,
	}
}

func (r *CampaignRepository) CreateCampaign(d *model.Campaign) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert campaign: %w", err)
	}

	return nil
}

func (r *CampaignRepository) CampaignResults(d *model.EmailCampaignResult) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert campaign result: %w", err)
	}
	return nil
}

func (r *CampaignRepository) GetScheduledCampaigns() {}

func (r *CampaignRepository) GetDraftCampaigns() {}

func (r *CampaignRepository) GetAllCampaigns() {}

func (r *CampaignRepository) GetSentCampaigns() {}
