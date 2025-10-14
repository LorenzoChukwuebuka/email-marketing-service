package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/campaigns/dto"
	"email-marketing-service/core/handler/campaigns/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	worker "email-marketing-service/internal/workers"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Service struct {
	store db.Store
	wrkr  *worker.Worker
}

func NewCampaignService(store db.Store, wrkr *worker.Worker) *Service {
	return &Service{
		store: store,
		wrkr:  wrkr,
	}
}

type CampaignWithGroupsResponse struct {
	*dto.CampaignResponseDTO
	Groups []*dto.GetCampaignContactGroupsResponse `json:"groups"`
}

func (s *Service) CreateCampaign(ctx context.Context, req *dto.CampaignDTO) (*dto.CampaignDTO, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
		"user":    req.UserId,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	exists, err := s.store.CheckCampaignNameExists(ctx, db.CheckCampaignNameExistsParams{
		CompanyID: _uuid["company"],
		Name:      req.Name,
	})
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if exists {
		return nil, common.ErrRecordExists
	}

	_, err = s.store.CreateCampaign(ctx, db.CreateCampaignParams{
		CompanyID:      _uuid["company"],
		Name:           req.Name,
		Subject:        sql.NullString{String: req.Subject, Valid: true},
		PreviewText:    sql.NullString{String: req.PreviewText, Valid: true},
		UserID:         _uuid["user"],
		SenderFromName: sql.NullString{String: req.SenderFromName, Valid: true},
		TemplateID:     uuid.NullUUID{UUID: _uuid["template"], Valid: true},
		RecipientInfo:  sql.NullString{String: req.RecipientInfo, Valid: true},
		Status:         sql.NullString{String: string(enums.CampaignStatus(enums.Draft)), Valid: true},
		TrackType:      sql.NullString{String: req.TrackType, Valid: true},
		Sender:         sql.NullString{String: req.Sender, Valid: true},
		ScheduledAt:    sql.NullTime{Time: req.ScheduledAt, Valid: true},
		HasCustomLogo:  sql.NullBool{Bool: req.HasCustomLogo, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}

	//notify user of campaign creation
	notificationTitle := fmt.Sprintf("You have successfully created campaign with name '%s' ", req.Name)
	additionalField := "campaign_creation"
	payload := worker.UserNotificationPayload{
		UserId:           _uuid["user"],
		NotifcationTitle: notificationTitle,
		AdditionalField:  &additionalField,
	}

	if taskId, err := s.wrkr.EnqueueTask(ctx, worker.TaskSendUserNotification, payload); err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	} else {
		log.Printf("Task enqueued with ID: %d", taskId)
	}

	return req, nil
}

func (s *Service) GetAllCampaigns(ctx context.Context, req *dto.FetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
		"user":    req.UserID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaigns, err := s.store.ListCampaignsByCompanyID(ctx, db.ListCampaignsByCompanyIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
		Column5:   req.SearchQuery,
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching record : %w", err)
	}

	count_campaigns, err := s.store.GetCampaignCounts(ctx, db.GetCampaignCountsParams{
		UserID:    _uuid["user"],
		CompanyID: _uuid["company"],
	})

	if err != nil {
		return nil, common.ErrFetchingCount
	}

	response := mapper.MapCampaignResponses(campaigns)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}
	data := common.Paginate(int(count_campaigns), items, req.Offset, req.Limit)
	return data, nil
}

func (s *Service) GetSingleCampaign(ctx context.Context, req *dto.FetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserID,
		"campaign": req.CampaignID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaign, err := s.store.GetCampaignByID(ctx, db.GetCampaignByIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		ID:        _uuid["campaign"],
	})
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	//get contact groups
	campaign_group, err := s.store.GetCampaignContactGroups(ctx, _uuid["campaign"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	// Map the campaign data to the response DTO
	groupData := mapper.MapCampaignGroups(campaign_group)
	campaignData := mapper.MapGetCampaignResponse(campaign)

	// Fetch template separately if template_id exists
	if campaign.TemplateID.Valid && campaign.TemplateID.UUID != uuid.Nil {
		template, err := s.store.GetTemplateByIDWithoutType(ctx, db.GetTemplateByIDWithoutTypeParams{
			ID:     campaign.TemplateID.UUID,
			UserID: uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		})

		if err != nil {
			if err == sql.ErrNoRows {
				// Template not found, set to nil
				campaignData.Template = nil
			} else {
				return nil, fmt.Errorf("error fetching template: %w", err)
			}
		} else {
			// Map the template data
			templateData := mapper.MapTemplateFromSeparateQuery(template)
			campaignData.Template = templateData
		}
	}

	return &CampaignWithGroupsResponse{
		CampaignResponseDTO: campaignData,
		Groups:              groupData,
	}, nil
}

func (s *Service) UpdateCampaign(ctx context.Context, req *dto.CampaignDTO, campaignId string) error {
	uuidMap := map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserId,
		"campaign": campaignId,
	}
	if req.TemplateId != "" {
		uuidMap["template"] = req.TemplateId
	}

	_uuid, err := common.ParseUUIDMap(uuidMap)
	if err != nil {
		return common.ErrInvalidUUID
	}

	_, err = s.store.UpdateCampaign(ctx, db.UpdateCampaignParams{
		ID:             _uuid["campaign"],
		UserID:         _uuid["user"],
		Name:           sql.NullString{String: req.Name, Valid: req.Name != ""},
		Subject:        sql.NullString{String: req.Subject, Valid: req.Subject != ""},
		PreviewText:    sql.NullString{String: req.PreviewText, Valid: req.PreviewText != ""},
		SenderFromName: sql.NullString{String: req.SenderFromName, Valid: req.SenderFromName != ""},
		TemplateID:     uuid.NullUUID{UUID: _uuid["template"], Valid: req.TemplateId != ""},
		RecipientInfo:  sql.NullString{String: req.RecipientInfo, Valid: req.RecipientInfo != ""},
		Status:         sql.NullString{String: string(req.Status), Valid: req.Status != ""},
		TrackType:      sql.NullString{String: req.TrackType, Valid: req.TrackType != ""},
		Sender:         sql.NullString{String: req.Sender, Valid: req.Sender != ""},
		ScheduledAt:    sql.NullTime{Time: req.ScheduledAt, Valid: req.ScheduledAt != time.Time{}},
		HasCustomLogo:  sql.NullBool{Bool: req.HasCustomLogo, Valid: true},
		IsPublished:    sql.NullBool{Bool: req.IsPublished, Valid: true},
		IsArchived:     sql.NullBool{Bool: req.IsArchived, Valid: true},
	})
	if err != nil {
		return err
	}

	notificationTitle := fmt.Sprintf("You have successfully updated campaign with name '%s' ", req.Name)
	additionalField := "campaign_update"
	payload := worker.UserNotificationPayload{
		UserId:           _uuid["user"],
		NotifcationTitle: notificationTitle,
		AdditionalField:  &additionalField,
	}

	if taskId, err := s.wrkr.EnqueueTask(ctx, worker.TaskSendUserNotification, payload); err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	} else {
		log.Printf("Task enqueued with ID: %d", taskId)
	}

	return nil
}

func (s *Service) DeleteCampaign(ctx context.Context, req *dto.FetchCampaignDTO) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyID,
		"user":     req.UserID,
		"campaign": req.CampaignID,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.SoftDeleteCampaign(ctx, db.SoftDeleteCampaignParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		ID:        _uuid["campaign"],
	})
	if err != nil {
		return common.ErrDeletingRecord
	}

	notificationTitle := fmt.Sprintf("You have successfully deleted campaign with id '%s' ", req.CampaignID)
	additionalField := "campaign_deletion"
	payload := &worker.UserNotificationPayload{
		UserId:           _uuid["user"],
		NotifcationTitle: notificationTitle,
		AdditionalField:  &additionalField,
	}

	if taskId, err := s.wrkr.EnqueueTask(ctx, worker.TaskSendUserNotification, payload); err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	} else {
		log.Printf("Task enqueued with ID: %d", taskId)
	}

	return nil
}

func (s *Service) CreateCampaignGroup(ctx context.Context, req *dto.CampaignGroupDTO) (*dto.CampaignGroupDTO, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": req.CampaignId,
		"group":    req.GroupId,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	//check if campaigngroup already exists with campaignid
	campaignGroupExists, err := s.store.CampaignGroupExists(ctx, _uuid["campaign"])
	if err != nil {
		return nil, fmt.Errorf("error checking if campaigngroup exists:%v", err)
	}

	if !campaignGroupExists {
		err := s.store.UpdateCampaignGroup(ctx, db.UpdateCampaignGroupParams{
			CampaignID:     _uuid["campaign"],
			ContactGroupID: _uuid["group"],
		})
		if err != nil {
			return nil, common.ErrUpdatingRecord
		}
	}

	err = s.store.CreateCampaignGroups(ctx, db.CreateCampaignGroupsParams{
		CampaignID:     _uuid["campaign"],
		ContactGroupID: _uuid["group"],
	})
	if err != nil {
		return nil, common.ErrCreatingRecord
	}

	return req, nil
}

func (s *Service) GetAllScheduledCampaigns(ctx context.Context, req *dto.FetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
		"user":    req.UserID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaigns, err := s.store.ListScheduledCampaignsByCompanyID(ctx, db.ListScheduledCampaignsByCompanyIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
		Column5:   req.SearchQuery,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching scheduled campaigns: %v", err)
	}

	count_campaigns, err := s.store.GetCampaignCounts(ctx, db.GetCampaignCountsParams{
		UserID:    _uuid["user"],
		CompanyID: _uuid["company"],
	})

	if err != nil {
		return nil, common.ErrFetchingCount
	}

	response := mapper.MapScheduledCampaignResponses(campaigns)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}
	data := common.Paginate(int(count_campaigns), items, req.Offset, req.Limit)
	return data, nil
}

func (s *Service) SendCampaign(ctx context.Context, req *dto.SendCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company":  req.CompanyId,
		"user":     req.UserId,
		"campaign": req.CampaignId,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	campaign, err := s.store.GetCampaignByID(ctx, db.GetCampaignByIDParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		ID:        _uuid["campaign"],
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if campaign.SentAt.Valid && !campaign.SentAt.Time.IsZero() {
		// Check if the campaign status is "sent" - if so, don't allow resending
		if campaign.Status.Valid && campaign.Status.String == string(enums.CampaignStatus(enums.Sent)) {
			return nil, fmt.Errorf("campaign already sent successfully")
		}
		// If status is "failed", we allow resending
		if campaign.Status.Valid && campaign.Status.String == string(enums.CampaignStatus(enums.Failed)) {
			log.Printf("Retrying failed campaign: %s", campaign.ID)
		}
	}

	if campaign.IsArchived.Valid && campaign.IsArchived.Bool {
		return nil, fmt.Errorf("campaign is archived")
	}

	// Check if the campaign is scheduled and not due yet
	if campaign.ScheduledAt.Valid && !campaign.ScheduledAt.Time.IsZero() {
		scheduledTime := campaign.ScheduledAt.Time

		if scheduledTime.After(time.Now()) {
			return nil, nil // Not due yet, exit without sending
		}
	}

	err = s.store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
		Status: sql.NullString{String: string(enums.CampaignStatus(enums.Queued)), Valid: true},
		ID:     _uuid["campaign"],
		UserID: _uuid["user"],
	})

	if err != nil {
		return nil, fmt.Errorf("error updating campaign status: %v", err)
	}

	payload := worker.SendCampaignEmailPayload{
		CompanyID:  _uuid["company"],
		UserID:     _uuid["user"],
		CampaignID: _uuid["campaign"],
	}

	taskID, err := s.wrkr.EnqueueTask(ctx, worker.TaskSendCampaignEmail, payload)
	if err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	} else {
		log.Printf("Task enqueued with ID: %d", taskID)
	}

	//notification
	notificationTitle := fmt.Sprintf("Your campaign with name '%s'  is currently processing. We will notify you when it is done", campaign.Name)
	additionalField := "campaign_processing"
	payload__ := &worker.UserNotificationPayload{
		UserId:           _uuid["user"],
		NotifcationTitle: notificationTitle,
		AdditionalField:  &additionalField,
	}

	if taskId, err := s.wrkr.EnqueueTask(ctx, worker.TaskSendUserNotification, payload__); err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	} else {
		log.Printf("Task enqueued with ID: %d", taskId)
	}

	return nil, nil
}

func (s *Service) GetAllRecipientsForACampaign(ctx context.Context, campaignId string, companyId string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
		"company":  companyId,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	result, err := s.store.GetEmailCampaignResultsByCampaign(ctx, db.GetEmailCampaignResultsByCampaignParams{
		CampaignID: _uuid["campaign"],
		CompanyID:  _uuid["company"],
	})
	if err != nil {
		return nil, err
	}

	data := mapper.MapCampaignEmailResponse(result)
	return data, nil
}

func (s *Service) GetEmailResultStats(ctx context.Context, campaignId string, companyId string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
		"company":  companyId,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	result, err := s.store.GetEmailCampaignStats(ctx, db.GetEmailCampaignStatsParams{
		CampaignID: _uuid["campaign"],
		CompanyID:  _uuid["company"],
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetUserCampaignStats(ctx context.Context, userID string) (map[string]int64, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	stats, err := s.store.GetUserCampaignStats(ctx, _uuid["user"])

	if err != nil {
		return nil, fmt.Errorf("error fetching user campaign stats: %w", err)
	}
	// Convert to int64 with proper type assertions
	totalEmailsSent := int64(stats.TotalEmailsSent)
	totalOpens := stats.TotalOpens.(int64)
	uniqueOpens := int64(stats.UniqueOpens)
	totalClicks := stats.TotalClicks.(int64)
	uniqueClicks := int64(stats.UniqueClicks)
	softBounces := int64(stats.SoftBounces)
	hardBounces := int64(stats.HardBounces)
	totalBounces := int64(stats.TotalBounces)
	totalDeliveries := int64(stats.TotalDeliveries)

	// Calculate open rate
	var openRate float64
	if totalEmailsSent > 0 {
		openRate = (float64(uniqueOpens) / float64(totalEmailsSent)) * 100
	}

	result := map[string]int64{
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

	return result, nil
}

func (s *Service) GetAllCampaignStatsByUser(ctx context.Context, req *dto.FetchCampaignDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": req.UserID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	// First get all campaigns for the user
	campaigns, err := s.store.GetAllCampaignsByUser(ctx, db.GetAllCampaignsByUserParams{
		UserID: _uuid["user"],
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching campaigns for user: %w", err)
	}

	//get campaign counts

	count_campaigns, err := s.store.GetCampaignCounts(ctx, db.GetCampaignCountsParams{
		UserID:    _uuid["user"],
		CompanyID: _uuid["company"],
	})

	if err != nil {
		return nil, common.ErrFetchingCount
	}

	var allCampaignStats []map[string]interface{}

	// For each campaign, get its individual stats
	for _, campaign := range campaigns {
		stats, err := s.store.GetCampaignStats(ctx, campaign.CampaignID)
		if err != nil {
			return nil, fmt.Errorf("error fetching stats for campaign %s: %w", campaign.CampaignID, err)
		}

		campaignStats := map[string]interface{}{
			"campaign_id":  campaign.CampaignID,
			"name":         campaign.Name,
			"recipients":   stats.TotalEmailsSent,
			"opened":       stats.UniqueOpens,
			"clicked":      stats.UniqueClicks,
			"unsubscribed": stats.Unsubscribed,
			"complaints":   stats.Complaints,
			"bounces":      stats.TotalBounces,
			"sent_date":    campaign.SentAt.Time,
		}

		allCampaignStats = append(allCampaignStats, campaignStats)
	}

	items := make([]any, len(allCampaignStats))
	for i, v := range allCampaignStats {
		items[i] = v
	}

	data := common.Paginate(int(count_campaigns), items, req.Offset, req.Limit)

	return data, nil
}

func (s *Service) GetCampaignStats(ctx context.Context, campaignId string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	result, err := s.store.GetCampaignStats(ctx, _uuid["campaign"])
	if err != nil {
		return nil, err
	}

	return result, nil
}
