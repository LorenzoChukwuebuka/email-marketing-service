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
		SentTemplateId: data.SentTemplateId,
		RecipientInfo:  data.RecipientInfo,
		TemplateId:     data.TemplateId,
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
	var campaign model.Campaign

	result := r.DB.Model(&campaign).
		Where("uuid = ? AND user_id = ?", campaignId, userId).
		Preload("CampaignGroups").
		Preload("Template", func(db *gorm.DB) *gorm.DB {
			return db.Where("uuid = templates.uuid").
				Select("id, uuid, template_name, sender_name, from_email, subject, type, email_html, email_design, is_editable, is_published, is_public_template, is_gallery_template, tags, description, image_url, is_active, editor_type")
		}).
		First(&campaign)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	response := r.createCampaignMapping(campaign)

	return response, nil
}

func (r *CampaignRepository) UpdateCampaign(d *model.Campaign) error {
	var existingCampaign model.Campaign

	if err := r.DB.Where("uuid = ? AND user_id = ?", d.UUID, d.UserId).First(&existingCampaign).Error; err != nil {
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

func (r *CampaignRepository) DeleteCampaign(campaignId string, userId string) error {
	result := r.DB.Where("uuid = ? AND user_id = ?  ", campaignId, userId).
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

func (r *CampaignRepository) CreateEmailCampaignResult(d *model.EmailCampaignResult) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert campaign result: %w", err)
	}
	return nil
}

func (r *CampaignRepository) UpdateEmailCampaignResult(d *model.EmailCampaignResult) error {
	var existingCampaign model.EmailCampaignResult

	if err := r.DB.Where("campaign_id = ? AND recipient_email = ?", d.CampaignID, d.RecipientEmail).First(&existingCampaign).Error; err != nil {
		return fmt.Errorf("failed to find campaign for update: %w", err)
	}

	// Update fields only if they are not nil
	if d.RecipientName != nil {
		existingCampaign.RecipientName = d.RecipientName
	}
	if d.OpenedAt != nil {
		existingCampaign.OpenedAt = d.OpenedAt
	}
	if d.ClickedAt != nil {
		existingCampaign.ClickedAt = d.ClickedAt
	}
	if d.ConversionAt != nil {
		existingCampaign.ConversionAt = d.ConversionAt
	}
	if d.UnsubscribeAt != nil {
		existingCampaign.UnsubscribeAt = d.UnsubscribeAt
	}
	if d.BounceStatus != "" {
		existingCampaign.BounceStatus = d.BounceStatus
	}
	if d.ComplaintStatus {
		existingCampaign.ComplaintStatus = d.ComplaintStatus
	}
	if d.DeviceType != "" {
		existingCampaign.DeviceType = d.DeviceType
	}
	if d.Location != "" {
		existingCampaign.Location = d.Location
	}
	if d.RetryCount > 0 {
		existingCampaign.RetryCount = d.RetryCount
	}
	if d.Notes != "" {
		existingCampaign.Notes = d.Notes
	}
	if d.Version != "" {
		existingCampaign.Version = d.Version
	}

	if err := r.DB.Save(&existingCampaign).Error; err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	return nil
}
