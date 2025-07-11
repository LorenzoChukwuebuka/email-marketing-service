package services

import (
	"bytes"
	"context"
	"database/sql"
	"email-marketing-service/core/handler/campaigns/dto"
	"email-marketing-service/core/handler/campaigns/mapper"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/enums"
	smtpfactory "email-marketing-service/internal/factory/smtpFactory"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"
	//"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type Service struct {
	store db.Store
}

func NewCampaignService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

const BATCH_SIZE = 100

var (
	cfg = config.LoadEnv()
)

// Custom error types for better error handling
type CampaignError struct {
	Type    string
	Message string
	Err     error
}

func (e CampaignError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type BatchError struct {
	BatchIndex int
	Recipients []string
	Errors     []RecipientError
}

type RecipientError struct {
	Email string
	Error string
}

type SendCampaignUUIDStruct struct {
	CompanyID  uuid.UUID `json:"company_id"`
	UserID     uuid.UUID `json:"user_id"`
	CampaignID uuid.UUID `json:"campaign_id"`
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
	campaignData := mapper.MapCampaignResponse(db.ListCampaignsByCompanyIDRow(campaign))

	// Fetch template separately if template_id exists
	if campaign.TemplateID.Valid && campaign.TemplateID.UUID != uuid.Nil {
		template, err := s.store.GetTemplateByIDWithoutType(ctx, db.GetTemplateByIDWithoutTypeParams{
			ID:     campaign.TemplateID.UUID,
			UserID: _uuid["user"],
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
		return nil, fmt.Errorf("campaign already sent")
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

	//add the uuids to the struct for literally transport
	requuid := &SendCampaignUUIDStruct{
		CompanyID:  _uuid["company"],
		UserID:     _uuid["user"],
		CampaignID: _uuid["campaign"],
	}

	//Start the sending process in a goroutine
	go func() {
		// Create a new context for the goroutine to avoid cancellation issues
		bgCtx := context.Background()

		if err := s.processCampaign(bgCtx, requuid); err != nil {
			// Log the detailed error
			log.Printf("Campaign processing failed with detailed error: %v", err)

			// Update campaign status to failed
			if updateErr := s.store.UpdateCampaignStatus(bgCtx, db.UpdateCampaignStatusParams{
				Status: sql.NullString{String: string(enums.CampaignStatus(enums.Failed)), Valid: true},
				ID:     _uuid["campaign"],
				UserID: _uuid["user"],
			}); updateErr != nil {
				log.Printf("error occurred while updating status to failed: %v", updateErr)
			}

			// You could also store the error details in database for later analysis
			s.storeCampaignError(bgCtx, _uuid["campaign"], err)
		} else {
			// Update campaign status to sent
			if updateErr := s.store.UpdateCampaignStatus(bgCtx, db.UpdateCampaignStatusParams{
				Status: sql.NullString{String: string(enums.CampaignStatus(enums.Sent)), Valid: true},
				ID:     _uuid["campaign"],
				UserID: _uuid["user"],
			}); updateErr != nil {
				log.Printf("campaign sent successfully but failed to update status: %v", updateErr)
			}
		}
	}()

	return nil, nil
}

func (s *Service) processCampaign(ctx context.Context, d *SendCampaignUUIDStruct) error {
	//get all associated emails for this campaign (from the contact groups)
	contactEmails, err := s.store.GetCampaignContactEmails(ctx, d.CampaignID)
	if err != nil {
		log.Printf("error occurred while fetching contacts: %v", err)
		return CampaignError{
			Type:    "CONTACT_FETCH_ERROR",
			Message: "Failed to fetch campaign contacts",
			Err:     err,
		}
	}

	if len(contactEmails) == 0 {
		return CampaignError{
			Type:    "NO_CONTACTS_ERROR",
			Message: "No contacts found for campaign",
			Err:     nil,
		}
	}

	mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, d.CompanyID)
	if err != nil {
		return CampaignError{
			Type:    "USAGE_CHECK_ERROR",
			Message: "Failed to fetch email usage record",
			Err:     err,
		}
	}

	if mailUsageRecord.RemainingEmails.Int32 == 0 {
		return CampaignError{
			Type:    "QUOTA_EXCEEDED_ERROR",
			Message: "Email quota exceeded",
			Err:     nil,
		}
	}

	// Create a channel to collect ALL errors from goroutines
	errChan := make(chan BatchError, (len(contactEmails)+BATCH_SIZE-1)/BATCH_SIZE)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a mutex to synchronize mail usage updates
	var mu sync.Mutex

	// Track batch processing
	var batchErrors []BatchError
	totalBatches := 0

	//send emails in batches
	for i := 0; i < len(contactEmails); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(contactEmails) {
			end = len(contactEmails)
		}
		batch := contactEmails[i:end]
		batchIndex := totalBatches

		wg.Add(1)
		go func(ctx context.Context, batch []string, batchIdx int) {
			defer wg.Done()

			batchErr := s.sendEmailBatch(&mu, ctx, d, batch, batchIdx)
			if len(batchErr.Errors) > 0 {
				errChan <- batchErr
			}
		}(ctx, batch, batchIndex)

		totalBatches++
		// small delay between batches to avoid overwhelming the email server
		time.Sleep(5 * time.Second)
	}

	// Start a goroutine to close the error channel when all workers are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect ALL errors from all batches
	for batchErr := range errChan {
		log.Printf("Batch %d had %d errors out of %d recipients",
			batchErr.BatchIndex, len(batchErr.Errors), len(batchErr.Recipients))
		for _, recErr := range batchErr.Errors {
			log.Printf("  - %s: %s", recErr.Email, recErr.Error)
		}
		batchErrors = append(batchErrors, batchErr)
	}

	// Update campaign status (we do this even if there were some errors)
	if err = s.store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
		Status: sql.NullString{String: string(enums.CampaignStatus(enums.Sent)), Valid: true},
		ID:     d.CampaignID,
		UserID: d.UserID,
		SentAt: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		log.Printf("error occurred while updating campaign status: %v", err)
		return CampaignError{
			Type:    "STATUS_UPDATE_ERROR",
			Message: "Failed to update campaign status after processing",
			Err:     err,
		}
	}

	// If there were errors in batches, return detailed information
	if len(batchErrors) > 0 {
		totalErrors := 0
		for _, batchErr := range batchErrors {
			totalErrors += len(batchErr.Errors)
		}

		return CampaignError{
			Type: "BATCH_PROCESSING_ERRORS",
			Message: fmt.Sprintf("Campaign completed with %d errors across %d batches out of %d total batches. Total contacts: %d, Failed: %d",
				totalErrors, len(batchErrors), totalBatches, len(contactEmails), totalErrors),
			Err: fmt.Errorf("batch errors: %+v", batchErrors),
		}
	}
	log.Printf("Campaign processed successfully: %d batches, %d total contacts", totalBatches, len(contactEmails))
	return nil
}

func (s *Service) sendEmailBatch(mu *sync.Mutex, ctx context.Context, d *SendCampaignUUIDStruct, recipients []string, batchIndex int) BatchError {
	batchError := BatchError{
		BatchIndex: batchIndex,
		Recipients: recipients,
		Errors:     []RecipientError{},
	}

	//first thing fetch the campaign and extract the template id
	campaign, err := s.store.GetCampaignByID(ctx, db.GetCampaignByIDParams{
		CompanyID: d.CompanyID,
		UserID:    d.UserID,
		ID:        d.CampaignID,
	})

	if err != nil {
		// If we can't get campaign info, all recipients in this batch fail
		for _, recipient := range recipients {
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: fmt.Sprintf("Failed to fetch campaign: %v", err),
			})
		}
		return batchError
	}

	if !campaign.TemplateID.Valid || campaign.TemplateID.UUID == uuid.Nil {
		log.Printf("invalid template Id: %v", campaign.TemplateID)
		for _, recipient := range recipients {
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: "Invalid template ID: template ID is required",
			})
		}
		return batchError
	}

	if !campaign.TemplateEmailHtml.Valid {
		log.Printf("empty template exiting sending....")
		for _, recipient := range recipients {
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: "Template design is empty",
			})
		}
		return batchError
	}

	user, err := s.store.GetUserByID(ctx, d.UserID)
	if err != nil {
		log.Printf("error fetching user: %v", err)
		for _, recipient := range recipients {
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: fmt.Sprintf("Error fetching user: %v", err),
			})
		}
		return batchError
	}

	// Track successful sends in this batch
	var successCount int32 = 0

	// Process each recipient
	for _, recipient := range recipients {
		// Modify the email template with tracking info
		modifiedTemplate, err := s.addTrackingToTemplate(campaign.TemplateEmailHtml.String, campaign.ID.String(), recipient, user.CompanyID.String())
		if err != nil {
			log.Printf("error adding tracking to template for recipient %s: %v", recipient, err)
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: fmt.Sprintf("Failed to add tracking to template: %v", err),
			})
			continue
		}

		// Send the email and track success with detailed error
		emailSent, sendErr := s.sendEmailWithError(
			ctx,
			recipient,
			modifiedTemplate,
			campaign.Subject.String,
			campaign.PreviewText.String,
			user.Email,
			user.Fullname,
			d.UserID,
			d.CompanyID,
		)

		if emailSent {
			successCount++
		} else {
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: fmt.Sprintf("Failed to send email: %v", sendErr),
			})
		}

		// Create the EmailCampaignResult for tracking
		if err := s.createEmailCampaignResult(campaign.ID, d.CompanyID, recipient); err != nil {
			// This is not a critical error, but we should log it
			log.Printf("Warning: Failed to create campaign result for %s: %v", recipient, err)
		}
	}

	// Only update the database once per batch with the total count
	if successCount > 0 {
		mu.Lock()
		defer mu.Unlock() // Ensure unlock happens even if there's an error

		// Get current usage record
		mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, d.CompanyID)
		if err != nil {
			log.Printf("error fetching mail usage record: %v", err)
			// Add this as an error for all successfully sent emails in this batch
			for _, recipient := range recipients {
				// Only add this error for recipients that were successfully sent
				found := false
				for _, existingErr := range batchError.Errors {
					if existingErr.Email == recipient {
						found = true
						break
					}
				}
				if !found {
					batchError.Errors = append(batchError.Errors, RecipientError{
						Email: recipient,
						Error: fmt.Sprintf("Email sent but failed to update usage record: %v", err),
					})
				}
			}
		} else {
			// Update with the count of successfully sent emails
			_, err = s.store.UpdateEmailsSentAndRemaining(ctx, db.UpdateEmailsSentAndRemainingParams{
				CompanyID:  d.CompanyID,
				EmailsSent: sql.NullInt32{Int32: successCount, Valid: true},
				ID:         mailUsageRecord.ID,
			})
			if err != nil {
				log.Printf("error updating mail usage: %v", err)
				// Add this as a warning for all successfully sent emails
				for _, recipient := range recipients {
					found := false
					for _, existingErr := range batchError.Errors {
						if existingErr.Email == recipient {
							found = true
							break
						}
					}
					if !found {
						batchError.Errors = append(batchError.Errors, RecipientError{
							Email: recipient,
							Error: fmt.Sprintf("Email sent but failed to update usage count: %v", err),
						})
					}
				}
			}
		}

		if len(batchError.Errors) == 0 {
			log.Printf("Successfully sent %d emails in batch %d", successCount, batchIndex)
		} else {
			log.Printf("Batch %d completed: %d sent, %d errors", batchIndex, successCount, len(batchError.Errors))
		}
	}

	return batchError
}

func (s *Service) createEmailCampaignResult(campaignId uuid.UUID, companyId uuid.UUID, recipientEmail string) error {
	_, err := s.store.CreateEmailCampaignResult(context.Background(), db.CreateEmailCampaignResultParams{
		CampaignID:     campaignId,
		CompanyID:      companyId,
		Version:        sql.NullString{String: "1", Valid: true},
		RecipientEmail: recipientEmail,
		SentAt:         sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return fmt.Errorf("error creating email campaign result: %w", err)
	}

	return nil
}

func (s *Service) addTrackingToTemplate(template string, campaignId string, recipientEmail string, companyId string) (string, error) {
	if template == "" {
		return "", fmt.Errorf("empty template provided")
	}

	// Create tracking pixel and unsubscribe link
	trackingPixel := fmt.Sprintf(`<img src="%s/misc/track/open/%s?email=%s" alt="" width="1" height="1" style="display:none;" />`, cfg.SERVER_URL, campaignId, url.QueryEscape(recipientEmail))
	unsubscribeLink := fmt.Sprintf(
		`<div style="margin-top: 20px; font-size: 12px; color: #666666; text-align: center;">
        <a href="%s/misc/unsubscribe?email=%s&campaign=%s&companyId=%s" target="_blank" style="color: #666666; text-decoration: underline;">Unsubscribe</a>
    </div>`, cfg.SERVER_URL, url.QueryEscape(recipientEmail), url.QueryEscape(campaignId), url.QueryEscape(companyId))

	// Inject tracking pixel and unsubscribe link at the end of the body or document
	if strings.Contains(template, "</body>") {
		template = strings.Replace(template, "</body>", trackingPixel+unsubscribeLink+"</body>", 1)
	} else {
		template += trackingPixel + unsubscribeLink
	}

	// Modify links for click tracking
	doc, err := html.Parse(strings.NewReader(template))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var modifyLinks func(*html.Node)
	modifyLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, a := range n.Attr {
				if a.Key == "href" && !strings.Contains(a.Val, "unsubscribe") {
					originalURL := a.Val
					trackingURL := fmt.Sprintf("%s/misc/track/click/%s?email=%s&url=%s", cfg.SERVER_URL, campaignId, recipientEmail, url.QueryEscape(originalURL))
					n.Attr[i].Val = trackingURL
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			modifyLinks(c)
		}
	}

	modifyLinks(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to render modified HTML: %w", err)
	}

	return buf.String(), nil
}

// Modified sendEmail to return detailed error information
func (s *Service) sendEmailWithError(ctx context.Context, recipient string, emailContent string, subject string, previewText string, from string, fromName string, userId uuid.UUID, companyId uuid.UUID) (bool, error) {
	validEmail := helper.IsValidEmail(recipient)

	if !validEmail {
		return false, fmt.Errorf("invalid email address format")
	}

	authUser, err := s.store.GetMasterSMTPKey(ctx, userId)
	if err != nil {
		log.Printf("error fetching master smtp key: %v", err)
		return false, fmt.Errorf("failed to fetch SMTP credentials: %w", err)
	}

	if authUser.Status != "active" {
		log.Printf("user status is not active")
		return false, fmt.Errorf("user status is not active")
	}

	authModel := &domain.SMTPAuthUser{
		Username: authUser.SmtpLogin,
		Password: authUser.Password,
	}
	sender := &domain.Sender{
		Email: from,
		Name:  &fromName,
	}

	receiver := domain.Recipient{
		Email: recipient,
	}

	request := &domain.EmailRequest{
		Sender:      *sender,
		To:          receiver,
		Subject:     subject,
		HtmlContent: &emailContent,
		PreviewText: &previewText,
		AuthUser:    *authModel,
	}

	emailBytes := []byte(emailContent)

	// Get the domain from the sender's email
	parts := strings.Split(from, "@")
	if len(parts) != 2 {
		log.Printf("invalid sender email format")
		return false, fmt.Errorf("invalid sender email format")
	}
	senderDomain := parts[1]

	sender_domain, err := s.store.FindDomainByNameAndCompany(ctx, db.FindDomainByNameAndCompanyParams{
		Domain:    senderDomain,
		CompanyID: companyId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			// If the domain is not found, proceed without signing
			//the mail will be sent from the app's domain
			// Sanitize fromName for email use
			emailLocalPart := strings.ReplaceAll(strings.ToLower(fromName), " ", ".")
			request.Sender.Email = fmt.Sprintf("%s@%s", emailLocalPart, cfg.DOMAIN)

			if sendErr := s.sendEmailWithSMTP(request); sendErr != nil {
				return false, fmt.Errorf("SMTP send failed: %w", sendErr)
			}
			return true, nil
		}
		log.Printf("failed to fetch domain: %v", err)
		return false, fmt.Errorf("failed to fetch domain configuration: %w", err)
	}

	if sender_domain != (db.Domain{}) && sender_domain.Verified.Valid && sender_domain.Verified.Bool {
		// Check if the sender's domain matches or is a subdomain of the DKIM signing domain
		if !strings.HasSuffix(senderDomain, sender_domain.Domain) {
			log.Printf("sender domain %s does not align with DKIM signing domain %s", senderDomain, sender_domain.Domain)
			return false, fmt.Errorf("sender domain %s does not align with DKIM signing domain %s", senderDomain, sender_domain.Domain)
		}

		signedBody, err := s.signEmail(from, companyId, emailBytes)
		if err != nil {
			// Log the error, but continue with unsigned email
			log.Printf("failed to sign email: %v", err)
		} else {
			emailBytes = signedBody
			request.HtmlContent = (*string)(unsafe.Pointer(&emailBytes))
		}
	}

	if err := s.sendEmailWithSMTP(request); err != nil {
		return false, fmt.Errorf("SMTP send failed: %w", err)
	}
	return true, nil
}

func (s *Service) sendEmailWithSMTP(request *domain.EmailRequest) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	analyzer, err := helper.NewContentAnalyzer("internal/config/config.json", logger)
	if err != nil {
		return fmt.Errorf("failed to create content analyzer: %w", err)
	}

	// Analyze the content before sending the email
	analysisResult, err := analyzer.AnalyzeContent(context.TODO(), *request.HtmlContent, nil)
	if err != nil {
		return fmt.Errorf("failed to analyze content: %w", err)
	}

	// Check if the content is flagged as spam or contains suspicious patterns
	if !analysisResult.IsSafe {
		logger.Warn("Email content flagged as unsafe",
			zap.Float64("spam_score", analysisResult.SpamScore),
			zap.String("message", analysisResult.Message),
			zap.Strings("suspicious_patterns", analysisResult.SuspiciousPatterns),
		)
		return fmt.Errorf("email content flagged as unsafe: %s", analysisResult.Message)
	}

	mailS, err := smtpfactory.MailFactory(cfg.MAIL_PROCESSOR)
	if err != nil {
		return fmt.Errorf("failed to create mail factory: %w", err)
	}

	if err := mailS.HandleSendMail(request); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Service) signEmail(domainEmail string, companyId uuid.UUID, emailBody []byte) ([]byte, error) {
	// Fetch the domain associated with the sender
	domain, err := s.store.FindDomainByNameAndCompany(context.Background(), db.FindDomainByNameAndCompanyParams{
		Domain:    domainEmail,
		CompanyID: companyId,
	})

	if err != nil || !domain.Verified.Valid {
		return nil, fmt.Errorf("domain not found or not verified")
	}

	helper.ValidatePrivateKey(domain.DkimPrivateKey.String)

	// DKIM signing process
	signedEmail, err := helper.SignEmail(&emailBody, domain.Domain, domain.DkimSelector.String, string(domain.DkimPrivateKey.String))
	if err != nil {
		return nil, fmt.Errorf("failed to sign email: %v", err)
	}

	return signedEmail, nil
}

// Helper function to store campaign errors for later analysis
func (s *Service) storeCampaignError(ctx context.Context, campaignID uuid.UUID, err error) {
	dbErr := s.store.CreateCampaignError(ctx, db.CreateCampaignErrorParams{
		CampaignID:   campaignID,
		ErrorType:    sql.NullString{String: fmt.Sprintf("%v", err), Valid: true},
		ErrorMessage: err.Error(),
	})
	if dbErr != nil {
		log.Printf("Failed to store campaign error: %v", dbErr)
	}
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
