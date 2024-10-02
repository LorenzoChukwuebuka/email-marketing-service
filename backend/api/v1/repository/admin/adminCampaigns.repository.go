package adminrepository

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AdminCampaignRepository struct {
	DB *gorm.DB
}

func NewAdminCampaignRepository(db *gorm.DB) *AdminCampaignRepository {
	return &AdminCampaignRepository{
		DB: db,
	}

}

func (r *AdminCampaignRepository) GetAllUserCampaigns(userId string, search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var campaigns []model.Campaign

	query := r.DB.Preload("Template").Preload("CampaignGroups").Where("user_id = ?", userId).Order("created_at DESC")

	if search != "" {
		query = query.Where("name ILIKE ? OR subject ILIKE ? OR uuid ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	pagination, err := repository.Paginate(query, params, &campaigns)
	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate campaigns: %w", err)
	}

	pagination.Data = r.campaignsToResponses(campaigns)

	return pagination, nil
}

func (r *AdminCampaignRepository) GetAUserCampaign(uuid string) (*model.CampaignResponse, error) {
	var campaign model.Campaign
	if err := r.DB.Preload("Template").Preload("CampaignGroups").Where("uuid = ?", uuid).First(&campaign).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("campaign not found")
		}
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	response := r.campaignToResponse(&campaign)
	return &response, nil
}

func (r *AdminCampaignRepository) DeleteUserCampaign(uuid string) error {
	result := r.DB.Where("uuid = ?", uuid).Delete(&model.Campaign{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete campaign: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("campaign not found")
	}
	return nil
}

func (r *AdminCampaignRepository) SuspendUserCampaign(uuid string) error {
	result := r.DB.Model(&model.Campaign{}).Where("uuid = ?", uuid).Update("status", model.Suspended)
	if result.Error != nil {
		return fmt.Errorf("failed to suspend campaign: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("campaign not found")
	}
	return nil
}

func (r *AdminCampaignRepository) campaignsToResponses(campaigns []model.Campaign) []model.CampaignResponse {
	responses := make([]model.CampaignResponse, len(campaigns))
	for i, campaign := range campaigns {
		responses[i] = r.campaignToResponse(&campaign)
	}
	return responses
}

func (r *AdminCampaignRepository) campaignToResponse(campaign *model.Campaign) model.CampaignResponse {
	return model.CampaignResponse{
		UUID:           campaign.UUID,
		Name:           campaign.Name,
		Subject:        campaign.Subject,
		PreviewText:    campaign.PreviewText,
		UserId:         campaign.UserId,
		SenderFromName: campaign.SenderFromName,
		TemplateId:     campaign.TemplateId,
		SentTemplateId: campaign.SentTemplateId,
		RecipientInfo:  campaign.RecipientInfo,
		IsPublished:    campaign.IsPublished,
		Status:         string(campaign.Status),
		TrackType:      string(campaign.TrackType),
		IsArchived:     campaign.IsArchived,
		SentAt:         campaign.SentAt,
		Sender:         campaign.Sender,
		ScheduledAt:    formatTime(campaign.ScheduledAt),
		CreatedAt:      campaign.CreatedAt.Format(time.RFC3339),
		//	UpdatedAt:      formatTime(&campaign.UpdatedAt),

		Template:       campaign.Template,
		CampaignGroups: r.campaignGroupsToResponses(campaign.CampaignGroups),
	}
}

func (r *AdminCampaignRepository) campaignGroupsToResponses(groups []model.CampaignGroup) []model.CampaignGroupResponse {
	responses := make([]model.CampaignGroupResponse, len(groups))
	for i, group := range groups {
		responses[i] = model.CampaignGroupResponse{
			UUID:       group.UUID,
			CampaignId: group.CampaignId,
			GroupId:    group.GroupId,
			CreatedAt:  group.CreatedAt.Format(time.RFC3339),
		}
	}
	return responses
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(time.RFC3339)
	return &formatted
}
