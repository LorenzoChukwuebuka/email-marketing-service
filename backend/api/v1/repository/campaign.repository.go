package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type CampaignRepository struct {
	DB *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{
		DB: db,
	}
}

func (r *CampaignRepository) createCampaignMapping(data model.Campaign) *model.CampaignResponse {
	result := model.CampaignResponse{
		ID:             data.ID,
		UUID:           data.UUID,
		Name:           data.Name,
		Subject:        data.Subject,
		PreviewText:    data.PreviewText,
		UserId:         data.UserId,
		SenderFromName: data.SenderFromName,
		TemplateId:     data.TemplateId,
		SentTemplateId: data.SentTemplateId,
		RecipientInfo:  data.RecipientInfo,
		IsPublished:    data.IsPublished,
		Status:         string(data.Status), // Convert CampaignStatus to string
		TrackType:      string(data.TrackType),
		IsArchived:     data.IsArchived,
		SentAt:         data.SentAt,
		ScheduledAt: func() *string {
			if data.ScheduledAt != nil {
				htime := data.ScheduledAt.Format(time.RFC3339)
				return &htime
			}
			return nil
		}(),

		Template:  data.Template,
		Sender:    data.Sender,
		CreatedAt: data.CreatedAt.String(),
		UpdatedAt: data.UpdatedAt.String(),
		DeletedAt: func() *string {
			var htime string

			if data.DeletedAt.Valid {
				htime = data.DeletedAt.Time.Format(time.RFC3339)
			}
			return &htime
		}(),
	}

	// Map CampaignGroups
	for _, campaignGroup := range data.CampaignGroups {
		result.CampaignGroups = append(result.CampaignGroups, r.createCampaignGroupMapping(campaignGroup))
	}

	return &result
}

func (r *CampaignRepository) createCampaignGroupMapping(data model.CampaignGroup) model.CampaignGroupResponse {
	return model.CampaignGroupResponse{
		ID:         data.ID,
		UUID:       data.UUID,
		CampaignId: data.CampaignId,
		GroupId:    data.GroupId,
		CreatedAt:  data.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  data.UpdatedAt.Format(time.RFC3339),
		DeletedAt: func() *string {
			if data.DeletedAt.Valid {
				deletedAt := data.DeletedAt.Time.Format(time.RFC3339)
				return &deletedAt
			}
			return nil
		}(),
	}
}

func (r *CampaignRepository) CreateCampaign(d *model.Campaign) (string, error) {
	if err := r.DB.Create(&d).Error; err != nil {
		return "", fmt.Errorf("failed to insert campaign: %w", err)
	}

	return d.UUID, nil
}

func (r *CampaignRepository) CampaignExists(d *model.Campaign) (bool, error) {
	result := r.DB.Where("name = ? AND user_id =?", d.Name, d.UserId).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil

}

func (r *CampaignRepository) GetAllCampaigns(userId string, params PaginationParams) (PaginatedResult, error) {
	var campaigns []model.Campaign

	query := r.DB.Model(&campaigns).Where("user_id = ?", userId).Preload("CampaignGroups").Order("created_at DESC")

	paginatedResult, err := Paginate(query, params, &campaigns)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	var response []model.CampaignResponse

	for _, campaign := range campaigns {
		response = append(response, *r.createCampaignMapping(campaign))
	}

	paginatedResult.Data = response

	return paginatedResult, nil
}

func (r *CampaignRepository) GetScheduledCampaigns(userId string, params PaginationParams) (PaginatedResult, error) {
	var campaigns []model.Campaign

	query := r.DB.Model(&campaigns).Where("user_id = ? AND status = ?", userId, model.Scheduled).Preload("CampaignGroups").Order("created_at DESC")

	paginatedResult, err := Paginate(query, params, &campaigns)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	var response []model.CampaignResponse

	for _, campaign := range campaigns {
		response = append(response, *r.createCampaignMapping(campaign))
	}

	paginatedResult.Data = response

	return paginatedResult, nil
}

func (r *CampaignRepository) GetSingleCampaign(userId string, campaignId string) (*model.CampaignResponse, error) {
	var campaigns model.Campaign

	result := r.DB.Model(&campaigns).Where("user_id = ?", userId).Preload("CampaignGroups").First(&campaigns)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	response := r.createCampaignMapping(campaigns)

	return response, nil
}

func (r *CampaignRepository) UpdateCampaign(d *model.Campaign) error {
	var existingCampaign model.Campaign

	if err := r.DB.Where("uuid = ?", d.UUID).First(&existingCampaign).Error; err != nil {
		return fmt.Errorf("failed to find campaign for update: %w", err)
	}

	if d.Name != "" {
		existingCampaign.Name = d.Name
	}
	if d.Subject != nil {
		existingCampaign.Subject = d.Subject
	}
	if d.PreviewText != nil {
		existingCampaign.PreviewText = d.PreviewText
	}

	if d.SenderFromName != nil {
		existingCampaign.SenderFromName = d.SenderFromName
	}
	if d.TemplateId != nil {
		existingCampaign.TemplateId = d.TemplateId
	}
	if d.SentTemplateId != nil {
		existingCampaign.SentTemplateId = d.SentTemplateId
	}
	if d.RecipientInfo != nil {
		existingCampaign.RecipientInfo = d.RecipientInfo
	}
	existingCampaign.IsPublished = d.IsPublished
	if d.Status != "" {
		existingCampaign.Status = d.Status
	}
	if d.TrackType != "" {
		existingCampaign.TrackType = d.TrackType
	}
	existingCampaign.IsArchived = d.IsArchived
	if d.SentAt != nil {
		existingCampaign.SentAt = d.SentAt
	}

	if d.Template != nil {
		existingCampaign.Template = d.Template
	}
	if d.Sender != nil {
		existingCampaign.Sender = d.Sender
	}
	if d.ScheduledAt != nil {
		existingCampaign.ScheduledAt = d.ScheduledAt
	}

	if err := r.DB.Save(&existingCampaign).Error; err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	return nil
}

func (r *CampaignRepository) DeleteCampaign(d *model.Campaign) error {
	result := r.DB.Where("uuid = ? AND user_id = ?  ", d.UUID, d.UserId).
		Delete(&model.Campaign{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove contact from group: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no matching record found to remove contact from group")
	}

	return nil
}

func (r *CampaignRepository) AddOrEditCampaignGroup(d *model.CampaignGroup) error {

	var existingGroup model.CampaignGroup
	if err := r.DB.Where("id = ?", d.ID).First(&existingGroup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			if err := r.DB.Create(d).Error; err != nil {
				return fmt.Errorf("error creating campaign group: %w", err)
			}
		} else {
			return fmt.Errorf("error fetching campaign group: %w", err)
		}
	} else {

		if err := r.DB.Model(&existingGroup).Updates(d).Error; err != nil {
			return fmt.Errorf("error updating campaign group: %w", err)
		}
	}

	return nil
}

func (r *CampaignRepository) CampaignResults(d *model.EmailCampaignResult) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert campaign result: %w", err)
	}
	return nil
}
