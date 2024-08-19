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
		SenderId:       data.SenderId,
		UserId:         data.UserId,
		SenderFromName: data.SenderFromName,
		TemplateId:     data.TemplateId,
		SentTemplateId: data.SentTemplateId,
		RecipientInfo:  data.RecipientInfo,
		IsPublished:    data.IsPublished,
		Status:         string(data.Status), // Convert CampaignStatus to string
		TrackType:      data.TrackType,
		IsArchived:     data.IsArchived,
		SentAt:         data.SentAt,
		CreatedBy:      data.CreatedBy,
		LastEditedBy:   data.LastEditedBy,
		Template:       data.Template,
		Sender:         data.Sender,
		CreatedAt:      data.CreatedAt.String(),
		UpdatedAt:      data.UpdatedAt.String(),
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

func (r *CampaignRepository) GetDraftCampaigns(userId string, params PaginationParams) (PaginatedResult, error) {
	var campaigns []model.Campaign

	query := r.DB.Model(&campaigns).Where("user_id = ? AND status = ?", userId, model.Draft).Preload("CampaignGroups").Order("created_at DESC")

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

func (r *CampaignRepository) GetSentCampaigns(userId string, params PaginationParams) (PaginatedResult, error) {
	var campaigns []model.Campaign

	query := r.DB.Model(&campaigns).Where("user_id = ? AND status = ?", userId, model.Sent).Preload("CampaignGroups").Order("created_at DESC")

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

func (r *CampaignRepository) UpdateCampaign(d model.Campaign) error {
	return r.DB.Model(&model.Campaign{}).Where("uuid = ?", d.UUID).Updates(d).Error
}

func (r *CampaignRepository) DeleteCampaign(d model.Campaign) error {
	return nil
}

func (r *CampaignRepository) CampaignResults(d *model.EmailCampaignResult) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert campaign result: %w", err)
	}
	return nil
}
