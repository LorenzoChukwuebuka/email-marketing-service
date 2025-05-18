package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
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

func (r *CampaignRepository) ConvertEmailCampaignResultToResponse(result *model.EmailCampaignResult) *model.EmailCampaignResultResponse {
	response := &model.EmailCampaignResultResponse{
		ID:             result.ID,
		CampaignID:     result.CampaignID,
		RecipientEmail: result.RecipientEmail,
		Version:        result.Version,
		SentAt:         result.SentAt,
		OpenedAt:       result.OpenedAt,
		OpenCount:      result.OpenCount,
		ClickedAt:      result.ClickedAt,
		ConversionAt:   result.ConversionAt,
		CreatedAt:      result.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      result.UpdatedAt.Format(time.RFC3339),
		DeletedAt:      nil,
	}

	if result.DeletedAt.Valid {
		deletedAt := result.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &deletedAt
	}

	return response
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

func (r *CampaignRepository) GetAllCampaigns(userId string, searchQuery string, params PaginationParams) (PaginatedResult, error) {
	var campaigns []model.Campaign

	query := r.DB.Model(&campaigns).Where("user_id = ?", userId).Preload("CampaignGroups")

	if searchQuery != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%")
	}

	query.Order("created_at DESC")

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
	var existingResult model.EmailCampaignResult
	result := r.DB.Where("recipient_email = ? AND campaign_id = ?", d.RecipientEmail, d.CampaignID).First(&existingResult)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing campaign result: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		return nil
	}

	// Create the new record
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert campaign result: %w", err)
	}

	return nil
}

func (r *CampaignRepository) GetEmailCampaignResult(campaignID, recipientEmail string) (*model.EmailCampaignResultResponse, error) {
	var emailCampaignResult model.EmailCampaignResult
	if err := r.DB.Where("campaign_id = ? AND recipient_email = ?", campaignID, recipientEmail).First(&emailCampaignResult).Error; err != nil {
		return nil, err
	}

	response := r.ConvertEmailCampaignResultToResponse(&emailCampaignResult)

	return response, nil

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

	if d.OpenCount > 0 {
		existingCampaign.OpenCount = d.OpenCount
	}

	if err := r.DB.Save(&existingCampaign).Error; err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	return nil
}

func (r *CampaignRepository) GetAllRecipientsForACampaign(campaignId string) (*[]model.EmailCampaignResultResponse, error) {
	var results []model.EmailCampaignResult
	var response []model.EmailCampaignResultResponse

	if err := r.DB.Where("campaign_id = ?", campaignId).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("error fetching campaign recipients: %w", err)
	}

	for _, result := range results {
		response = append(response, *r.ConvertEmailCampaignResultToResponse(&result))
	}

	return &response, nil
}

func (r *CampaignRepository) GetEmailResultStats(campaignID string) (map[string]interface{}, error) {
	var stats = make(map[string]interface{})

	// Total emails sent
	var totalSent int64
	if err := r.DB.Model(&model.EmailCampaignResult{}).Where("campaign_id = ?", campaignID).Count(&totalSent).Error; err != nil {
		return nil, fmt.Errorf("error calculating total emails sent: %w", err)
	}
	stats["total_emails_sent"] = totalSent

	// Unique opens
	var uniqueOpens int64
	if err := r.DB.Model(&model.EmailCampaignResult{}).Where("campaign_id = ? AND opened_at IS NOT NULL", campaignID).Count(&uniqueOpens).Error; err != nil {
		return nil, fmt.Errorf("error calculating unique opens: %w", err)
	}
	stats["unique_opens"] = uniqueOpens

	// Total opens
	var totalOpens struct {
		Total int64
	}
	if err := r.DB.Model(&model.EmailCampaignResult{}).Where("campaign_id = ?", campaignID).Select("COALESCE(SUM(open_count), 0) as total").Scan(&totalOpens).Error; err != nil {
		return nil, fmt.Errorf("error calculating total opens: %w", err)
	}
	stats["total_opens"] = totalOpens.Total

	// Unique clicks
	var uniqueClicks int64
	if err := r.DB.Model(&model.EmailCampaignResult{}).Where("campaign_id = ? AND clicked_at IS NOT NULL", campaignID).Count(&uniqueClicks).Error; err != nil {
		return nil, fmt.Errorf("error calculating unique clicks: %w", err)
	}
	stats["unique_clicks"] = uniqueClicks

	// Total clicks
	var totalClicks struct {
		Total int64
	}
	if err := r.DB.Model(&model.EmailCampaignResult{}).Where("campaign_id = ?", campaignID).Select("COALESCE(SUM(click_count), 0) as total").Scan(&totalClicks).Error; err != nil {
		return nil, fmt.Errorf("error calculating total clicks: %w", err)
	}
	stats["total_clicks"] = totalClicks.Total

	// Total bounces
	var totalBounces int64
	if err := r.DB.Model(&model.EmailCampaignResult{}).Where("campaign_id = ? AND bounce_status != ''", campaignID).Count(&totalBounces).Error; err != nil {
		return nil, fmt.Errorf("error calculating total bounces: %w", err)
	}
	stats["total_bounces"] = totalBounces

	return stats, nil
}

func (r *CampaignRepository) GetUserCampaignStats(userID string) (map[string]int64, error) {
	// Initialize variables to store aggregated stats
	var totalEmailsSent, totalOpens, uniqueOpens, totalClicks, uniqueClicks, softBounces, hardBounces int64 = 0, 0, 0, 0, 0, 0, 0

	// Get all campaigns for the user
	var campaigns []model.Campaign
	if err := r.DB.Where("user_id = ?", userID).Find(&campaigns).Error; err != nil {
		return nil, fmt.Errorf("error fetching campaigns for user: %w", err)
	}

	// Check if the user has any campaigns
	if len(campaigns) == 0 {
		return map[string]int64{}, nil // Return an empty map if no campaigns found
	}

	// Extract campaign UUIDs for use in the next query
	campaignIDs := make([]string, len(campaigns))
	for i, campaign := range campaigns {
		campaignIDs[i] = campaign.UUID
	}

	// Get total number of emails sent across all campaigns for the user
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ?", campaignIDs).
		Count(&totalEmailsSent).Error; err != nil {
		return nil, fmt.Errorf("error counting total emails sent: %w", err)
	}

	// Get total number of opens
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ?", campaignIDs).
		Select("COALESCE(SUM(open_count), 0)").
		Scan(&totalOpens).Error; err != nil {
		return nil, fmt.Errorf("error calculating total opens: %w", err)
	}

	// Get unique opens count
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ? AND open_count > 0", campaignIDs).
		Count(&uniqueOpens).Error; err != nil {
		return nil, fmt.Errorf("error calculating unique opens: %w", err)
	}

	// Get total number of clicks
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ?", campaignIDs).
		Select("COALESCE(SUM(click_count), 0)").
		Scan(&totalClicks).Error; err != nil {
		return nil, fmt.Errorf("error calculating total clicks: %w", err)
	}

	// Get unique clicks count
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ? AND click_count > 0", campaignIDs).
		Count(&uniqueClicks).Error; err != nil {
		return nil, fmt.Errorf("error calculating unique clicks: %w", err)
	}

	// Get count of soft bounces
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ? AND bounce_status = ?", campaignIDs, "soft").
		Count(&softBounces).Error; err != nil {
		return nil, fmt.Errorf("error calculating soft bounces: %w", err)
	}

	// Get count of hard bounces
	if err := r.DB.Model(&model.EmailCampaignResult{}).
		Where("campaign_id IN ? AND bounce_status = ?", campaignIDs, "hard").
		Count(&hardBounces).Error; err != nil {
		return nil, fmt.Errorf("error calculating hard bounces: %w", err)
	}

	// Calculate total bounces and deliveries
	totalBounces := softBounces + hardBounces
	totalDeliveries := totalEmailsSent - totalBounces

	// Calculate the open rate percentage
	var openRate float64
	if totalEmailsSent > 0 {
		openRate = (float64(uniqueOpens) / float64(totalEmailsSent)) * 100
	}

	// Construct and return the result map
	stats := map[string]int64{
		"total_emails_sent": totalEmailsSent,
		"total_opens":       totalOpens,
		"unique_opens":      uniqueOpens,
		"total_clicks":      totalClicks,
		"unique_clicks":     uniqueClicks,
		"soft_bounces":      softBounces,
		"hard_bounces":      hardBounces,
		"total_bounces":     totalBounces,
		"total_deliveries":  totalDeliveries,
		"open_rate":         int64(openRate),
	}

	return stats, nil
}

func (r *CampaignRepository) GetAllCampaignStatsByUser(userID string) ([]map[string]interface{}, error) {
	var allCampaignStats []map[string]interface{}

	var campaigns []model.Campaign
	if err := r.DB.Where("user_id = ?", userID).Find(&campaigns).Error; err != nil {
		return nil, fmt.Errorf("error fetching campaigns for user: %w", err)
	}

	for _, campaign := range campaigns {
		stats, err := r.GetEmailResultStats(campaign.UUID)
		if err != nil {
			return nil, fmt.Errorf("error fetching stats for campaign %s: %w", campaign.UUID, err)
		}

		campaignStats := map[string]interface{}{
			"campaign_id":  campaign.UUID,
			"name":         campaign.Name,
			"recipients":   stats["total_emails_sent"],
			"opened":       stats["unique_opens"],
			"clicked":      stats["unique_clicks"],
			"unsubscribed": 0,
			"complaints":   0,
			"bounces":      stats["total_bounces"],
			"sent_date":    campaign.SentAt,
		}

		allCampaignStats = append(allCampaignStats, campaignStats)
	}

	return allCampaignStats, nil
}

func (r *CampaignRepository) GetDueScheduledCampaigns() ([]model.CampaignResponse, error) {
	var campaigns []model.Campaign
	currentTime := time.Now().Format(time.RFC3339)
	err := r.DB.Where("scheduled_at <= ? AND status = ?", currentTime, model.Scheduled).Find(&campaigns).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching due scheduled campaigns: %w", err)
	}

	var response []model.CampaignResponse
	for _, campaign := range campaigns {
		campaignResponse := r.createCampaignMapping(campaign)
		response = append(response, *campaignResponse)
	}

	return response, nil
}
