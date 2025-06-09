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
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"
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

type SendCampaignUUIDStruct struct {
	CompanyID  uuid.UUID `json:"company_id"`
	UserID     uuid.UUID `json:"user_id"`
	CampaignID uuid.UUID `json:"campaign_id"`
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

func (s *Service) GetSingleCampaign(ctx context.Context, req *dto.FetchCampaignDTO) (*dto.CampaignResponseDTO, error) {
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

	data := mapper.MapCampaignResponse(db.ListCampaignsByCompanyIDRow(campaign))
	return data, nil
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
		Name:           req.Name,
		Subject:        sql.NullString{String: req.Subject, Valid: true},
		PreviewText:    sql.NullString{String: req.PreviewText, Valid: true},
		SenderFromName: sql.NullString{String: req.SenderFromName, Valid: true},
		TemplateID:     uuid.NullUUID{UUID: _uuid["template"], Valid: true},
		RecipientInfo:  sql.NullString{String: req.RecipientInfo, Valid: true},
		Status:         sql.NullString{String: string(req.Status), Valid: true},
		TrackType:      sql.NullString{String: req.TrackType, Valid: true},
		Sender:         sql.NullString{String: req.Sender, Valid: true},
		ScheduledAt:    sql.NullTime{Time: req.ScheduledAt, Valid: true},
		HasCustomLogo:  sql.NullBool{Bool: req.HasCustomLogo, Valid: true},
		IsPublished:    sql.NullBool{Bool: req.IsPublished, Valid: true},
		IsArchived:     sql.NullBool{Bool: req.IsArchived, Valid: true},
	})
	if err != nil {
		return common.ErrUpdatingRecord
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
		return nil, common.ErrFetchingRecord
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
		return nil, common.ErrUpdatingRecord
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

		if s.processCampaign(bgCtx, requuid); err != nil {
			if err = s.store.UpdateCampaignStatus(bgCtx, db.UpdateCampaignStatusParams{
				Status: sql.NullString{String: string(enums.CampaignStatus(enums.Failed)), Valid: true},
				ID:     _uuid["campaign"],
				UserID: _uuid["user"],
			}); err != nil {
				log.Printf("error occured while updating status :%v", err)
			}

			log.Printf("error occurred while sending mails :%v", err)
		} else {
			if err = s.store.UpdateCampaignStatus(bgCtx, db.UpdateCampaignStatusParams{
				Status: sql.NullString{String: string(enums.CampaignStatus(enums.Sent)), Valid: true},
				ID:     _uuid["campaign"],
				UserID: _uuid["user"],
			}); err != nil {
				log.Printf("error occured while updating status :%v", err)
			}
		}
	}()

	return nil, nil
}

func (s *Service) processCampaign(ctx context.Context, d *SendCampaignUUIDStruct) error {
	//get all associated emails for this campaign (from the contact groups)
	contactEmails, err := s.store.GetCampaignContactEmails(ctx, d.CampaignID)
	if err != nil {
		log.Printf("error occured while fetching contacts :%v", err)
		return err
	}

	mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, d.CompanyID)
	if err != nil {
		return fmt.Errorf("error fetching or creating mail usage record: %w", err)
	}

	if mailUsageRecord.RemainingEmails.Int32 == 0 {
		return fmt.Errorf("you have exceeded your plan limit")
	}

	if err = s.store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
		Status: sql.NullString{String: string(enums.CampaignStatus(enums.Sent)), Valid: true},
		ID:     d.CampaignID,
		UserID: d.UserID,
	}); err != nil {
		log.Printf("error occured while updating status :%v", err)
	}

	// Create a channel to receive errors from goroutines
	errChan := make(chan error, len(contactEmails))

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a mutex to synchronize mail usage updates
	var mu sync.Mutex

	//send emails in batches
	for i := 0; i < len(contactEmails); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(contactEmails) {
			end = len(contactEmails)
		}
		batch := contactEmails[i:end]

		wg.Add(1)
		go func(ctx context.Context, batch []string) {
			defer wg.Done()
			err := s.sendEmailBatch(&mu, ctx, d, batch)
			if err != nil {
				errChan <- err
			}
		}(ctx, batch)
		// small delay between batches to avoid overwhelming the email server
		time.Sleep(5 * time.Second)
	}

	// Start a goroutine to close the error channel when all workers are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors from the goroutines
	for err := range errChan {
		return err
	}

	return nil
}

func (s *Service) sendEmailBatch(mu *sync.Mutex, ctx context.Context, d *SendCampaignUUIDStruct, recipients []string) error {
	//first thing fetch the campaign and extract the template id
	campaign, err := s.store.GetCampaignByID(ctx, db.GetCampaignByIDParams{
		CompanyID: d.CompanyID,
		UserID:    d.UserID,
		ID:        d.CampaignID,
	})

	if err != nil {
		return common.ErrFetchingRecord
	}

	if !campaign.TemplateID.Valid || campaign.TemplateID.UUID == uuid.Nil {
		log.Printf("invalid template Id :%v", campaign.TemplateID)
		return fmt.Errorf("invalid template ID: template ID is required")
	}

	if !campaign.TemplateEmailHtml.Valid {
		log.Printf("empty template exiting sending....")
		return fmt.Errorf("template design is empty")
	}

	user, err := s.store.GetUserByID(ctx, d.UserID)
	if err != nil {
		log.Printf("error fetching user:%v", err)
		return fmt.Errorf("error fetching user:%w", err)
	}

	// Track successful sends in this batch
	var successCount int32 = 0

	// Process each recipient
	for _, recipient := range recipients {
		// Modify the email template with tracking info
		modifiedTemplate, err := s.addTrackingToTemplate(campaign.TemplateEmailHtml.String, campaign.ID.String(), recipient, user.CompanyID.String())
		if err != nil {
			log.Printf("error adding tracking to template for recipient %s: %v", recipient, err)
			continue // Skip this recipient but continue with others
		}

		// Send the email and track success
		emailSent := s.sendEmail(
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
		}

		// Create the EmailCampaignResult for tracking
		err = s.createEmailCampaignResult(campaign.ID, recipient)
		if err != nil {
			return fmt.Errorf("error creating email campaign result for recipient %s: %w", recipient, err)
		}

	}

	// Only update the database once per batch with the total count
	if successCount > 0 {
		mu.Lock()
		defer mu.Unlock() // Ensure unlock happens even if there's an error

		// Get current usage record
		mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, d.CompanyID)
		if err != nil {
			return fmt.Errorf("error fetching mail usage record: %w", err)
		}

		// Update with the count of successfully sent emails
		_, err = s.store.UpdateEmailsSentAndRemaining(ctx, db.UpdateEmailsSentAndRemainingParams{
			CompanyID:  d.CompanyID,
			EmailsSent: sql.NullInt32{Int32: successCount, Valid: true},
			ID:         mailUsageRecord.ID,
		})
		if err != nil {
			return fmt.Errorf("error updating mail usage: %w", err)
		}

		log.Printf("Successfully sent %d emails in this batch", successCount)
	}

	return nil
}

func (s *Service) createEmailCampaignResult(campaignId uuid.UUID, recipientEmail string) error {
	_, err := s.store.CreateEmailCampaignResult(context.Background(), db.CreateEmailCampaignResultParams{
		CampaignID:     campaignId,
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
	trackingPixel := fmt.Sprintf(`<img src="%s/campaigns/track/open/%s?email=%s" alt="" width="1" height="1" style="display:none;" />`, cfg.SERVER_URL, campaignId, url.QueryEscape(recipientEmail))
	unsubscribeLink := fmt.Sprintf(
		`<div style="margin-top: 20px; font-size: 12px; color: #666666; text-align: center;">
        <a href="%s/campaigns/unsubscribe?email=%s&campaign=%s&companyId=%s" target="_blank" style="color: #666666; text-decoration: underline;">Unsubscribe</a>
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
					trackingURL := fmt.Sprintf("%s/campaigns/track/click/%s?email=%s&url=%s", cfg.SERVER_URL, campaignId, recipientEmail, url.QueryEscape(originalURL))
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

// Helper function to simulate/handle actual email sending
func (s *Service) sendEmail(ctx context.Context, recipient string, emailContent string, subject string, previewText string, from string, fromName string, userId uuid.UUID, companyId uuid.UUID) bool {
	validEmail := helper.IsValidEmail(recipient)

	if !validEmail {
		return false
	}

	authUser, err := s.store.GetMasterSMTPKey(ctx, userId)
	if err != nil {
		log.Printf("error fetching master smtp key: %v", err)
		return false
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
		return false
	}
	senderDomain := parts[1]

	sender_domain, err := s.store.FindDomainByNameAndCompany(ctx, db.FindDomainByNameAndCompanyParams{
		Domain:    senderDomain,
		CompanyID: companyId,
	})

	if err != nil {
		if err != sql.ErrNoRows {
			// If the domain is not found, proceed without signing
			s.sendEmailWithSMTP(request)
			return true
		}
		log.Printf("failed to fetch domain: %v", err)
		return false
	}

	if sender_domain != (db.Domain{}) && sender_domain.Verified.Valid && sender_domain.Verified.Bool {
		// Check if the sender's domain matches or is a subdomain of the DKIM signing domain
		if !strings.HasSuffix(senderDomain, sender_domain.Domain) {
			log.Printf("sender domain %s does not align with DKIM signing domain %s", senderDomain, sender_domain.Domain)
			return false
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
	s.sendEmailWithSMTP(request)
	return true
}

func (s *Service) sendEmailWithSMTP(request *domain.EmailRequest) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	analyzer, err := helper.NewContentAnalyzer("config.json", logger)
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

	// Load the DKIM private key
	// privateKey, err := base64.StdEncoding.DecodeString(domain.DKIMPrivateKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to decode DKIM private key: %v", err)
	// }

	helper.ValidatePrivateKey(domain.DkimPrivateKey.String)

	// DKIM signing process
	signedEmail, err := helper.SignEmail(&emailBody, domain.Domain, domain.DkimSelector.String, string(domain.DkimPrivateKey.String))
	if err != nil {
		return nil, fmt.Errorf("failed to sign email: %v", err)
	}

	return signedEmail, nil
}

// func (s *CampaignService) GetAllRecipientsForACampaign(campaignId string) (*[]model.EmailCampaignResultResponse, error) {
// 	result, err := s.CampaignRepo.GetAllRecipientsForACampaign(campaignId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (s *CampaignService) GetEmailResultStats(campaignId string) (map[string]interface{}, error) {
// 	result, err := s.CampaignRepo.GetEmailResultStats(campaignId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (s *CampaignService) GetUserCampaignStats(userId string) (map[string]int64, error) {
// 	userStats, err := s.CampaignRepo.GetUserCampaignStats(userId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return userStats, nil
// }

// func (s *CampaignService) GetUserCampaignsStats(userId string) ([]map[string]interface{}, error) {
// 	userStats, err := s.CampaignRepo.GetAllCampaignStatsByUser(userId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return userStats, nil
// }
