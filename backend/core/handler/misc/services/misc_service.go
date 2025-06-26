package services

import (
	"context"
	"database/sql"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	smtpfactory "email-marketing-service/internal/factory/smtpFactory"
	"email-marketing-service/internal/helper"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unsafe"
	//"github.com/pkg/errors"
	"email-marketing-service/core/handler/admin/plans/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MiscService struct {
	store db.Store
}

func NewMiscService(store db.Store) *MiscService {
	return &MiscService{
		store: store,
	}

}

const BATCH_SIZE = 100

type BatchError struct {
	BatchIndex int
	Recipients []string
	Errors     []RecipientError
}

type RecipientError struct {
	Email string
	Error string
}

var (
	cfg = config.LoadEnv()
)

func (s *MiscService) TrackOpenCampaignEmails(ctx context.Context, campaignId string, email string, deviceType string, ipAddress string) error {
	htime := time.Now().UTC()

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
	})

	if err != nil {
		return common.ErrInvalidUUID
	}

	fmt.Printf("campaign:%v,email:%s", _uuid["campaign"], email)

	// Retrieve the existing email campaign result
	existingEmailResult, err := s.store.GetEmailCampaignResult(ctx, db.GetEmailCampaignResultParams{CampaignID: _uuid["campaign"], RecipientEmail: email})
	if err != nil {
		return err
	}

	// Increment the OpenCount by 1
	var openCount int32 = 1

	if existingEmailResult != (db.EmailCampaignResult{}) {
		openCount = existingEmailResult.OpenCount.Int32 + 1
	}

	_, err = s.store.UpdateEmailCampaignResult(ctx, db.UpdateEmailCampaignResultParams{
		CampaignID:     _uuid["campaign"],
		RecipientEmail: email,
		OpenedAt:       sql.NullTime{Time: htime, Valid: true},
		Location:       sql.NullString{String: ipAddress, Valid: true},
		DeviceType:     sql.NullString{String: deviceType, Valid: true},
		OpenCount:      sql.NullInt32{Int32: openCount, Valid: true},
	})

	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}

func (s *MiscService) UnsubscribeFromCampaign(ctx context.Context, campaignId string, email string, companyId string) error {
	htime := time.Now().UTC()
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
		"company":  companyId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}
	err = s.store.ExecTx(ctx, func(q *db.Queries) error {

		_, err = q.UpdateEmailCampaignResult(ctx, db.UpdateEmailCampaignResultParams{
			CampaignID:     _uuid["campaign"],
			RecipientEmail: email,
			UnsubscribedAt: sql.NullTime{Time: htime, Valid: true},
		})

		if err != nil {
			return err
		}

		err = q.UpdateContactSubscription(ctx, db.UpdateContactSubscriptionParams{
			CompanyID:    _uuid["company"],
			Email:        email,
			IsSubscribed: sql.NullBool{Bool: false, Valid: true},
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}

func (s *MiscService) TrackClickedCampaignsEmails(ctx context.Context, campaignId string, email string) error {
	htime := time.Now().UTC()
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"campaign": campaignId,
	})

	if err != nil {
		return common.ErrInvalidUUID
	}

	// Retrieve the existing email campaign result
	existingEmailResult, err := s.store.GetEmailCampaignResult(ctx, db.GetEmailCampaignResultParams{CampaignID: _uuid["campaign"], RecipientEmail: email})
	if err != nil {
		return err
	}

	var clickCount int32 = 1
	if existingEmailResult != (db.EmailCampaignResult{}) {
		clickCount = existingEmailResult.OpenCount.Int32 + 1
	}

	_, err = s.store.UpdateEmailCampaignResult(ctx, db.UpdateEmailCampaignResultParams{
		CampaignID:     _uuid["campaign"],
		RecipientEmail: email,
		ClickedAt:      sql.NullTime{Time: htime, Valid: true},
		ClickCount:     sql.NullInt32{Int32: clickCount, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *MiscService) GetPlans() (any, error) {
	return nil, nil
}

func (s *MiscService) PrepareMail(ctx context.Context, req *domain.EmailRequest, companyId uuid.UUID, userId uuid.UUID) (map[string]interface{}, error) {
	mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, companyId)
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if mailUsageRecord.RemainingEmails.Int32 == 0 {
		return nil, fmt.Errorf("email quota exceeded")
	}

	authUser, err := s.store.GetMasterSMTPKey(ctx, userId)
	if err != nil {
		log.Printf("error fetching master smtp key: %v", err)
		return nil, fmt.Errorf("failed to fetch SMTP credentials: %w", err)
	}

	if authUser.Status != "active" {
		log.Printf("user status is not active")
		return nil, fmt.Errorf("user status is not active")
	}

	authModel := &domain.SMTPAuthUser{
		Username: authUser.SmtpLogin,
		Password: authUser.Password,
	}

	// Check type of req.To and convert it to a slice of Recipient if necessary
	var recipients []domain.Recipient
	switch to := req.To.(type) {
	case domain.Recipient:
		recipients = []domain.Recipient{to}
	case []domain.Recipient:
		recipients = to
	default:
		return nil, fmt.Errorf("invalid recipient type")
	}

	// Create a channel to collect errors from goroutines
	errChan := make(chan BatchError, (len(recipients)+BATCH_SIZE-1)/BATCH_SIZE)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a mutex to synchronize mail usage updates
	var mu sync.Mutex

	// Track batch processing
	var batchErrors []BatchError
	totalBatches := 0

	// Process emails in batches
	for i := 0; i < len(recipients); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(recipients) {
			end = len(recipients)
		}
		batch := recipients[i:end]
		batchIndex := totalBatches

		wg.Add(1)
		go func(ctx context.Context, batch []domain.Recipient, batchIdx int) {
			defer wg.Done()

			batchErr := s.sendTransactionalEmailBatch(&mu, ctx, &TransactionalEmailParams{
				CompanyID:   companyId,
				UserID:      userId,
				AuthUser:    *authModel,
				Sender:      req.Sender,
				Subject:     req.Subject,
				HtmlContent: req.HtmlContent,
				Text:        req.Text,
			}, batch, batchIdx)

			if len(batchErr.Errors) > 0 {
				errChan <- batchErr
			}
		}(ctx, batch, batchIndex)

		totalBatches++
		// Small delay between batches to avoid overwhelming the email server
		time.Sleep(2 * time.Second) // Shorter delay for transactional emails
	}

	// Start a goroutine to close the error channel when all workers are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors from all batches
	for batchErr := range errChan {
		log.Printf("Transactional email batch %d had %d errors out of %d recipients",
			batchErr.BatchIndex, len(batchErr.Errors), len(batchErr.Recipients))
		for _, recErr := range batchErr.Errors {
			log.Printf("  - %s: %s", recErr.Email, recErr.Error)
		}
		batchErrors = append(batchErrors, batchErr)
	}

	// Return summary information
	result := map[string]interface{}{
		"total_batches":    totalBatches,
		"total_recipients": len(recipients),
		"errors":           len(batchErrors),
	}

	if len(batchErrors) > 0 {
		totalErrors := 0
		for _, batchErr := range batchErrors {
			totalErrors += len(batchErr.Errors)
		}
		result["failed_count"] = totalErrors
		result["success_count"] = len(recipients) - totalErrors
	} else {
		result["success_count"] = len(recipients)
		result["failed_count"] = 0
	}

	return result, nil
}

// TransactionalEmailParams holds parameters for sending transactional emails
type TransactionalEmailParams struct {
	CompanyID   uuid.UUID
	UserID      uuid.UUID
	AuthUser    domain.SMTPAuthUser
	Sender      domain.Sender
	Subject     string
	HtmlContent *string
	Text        *string
}

func (s *MiscService) sendTransactionalEmailBatch(mu *sync.Mutex, ctx context.Context, params *TransactionalEmailParams, recipients []domain.Recipient, batchIndex int) BatchError {
	batchError := BatchError{
		BatchIndex: batchIndex,
		Recipients: make([]string, len(recipients)), // Convert to string slice for consistency
		Errors:     []RecipientError{},
	}

	// Convert recipients to string slice for tracking
	for i, recipient := range recipients {
		batchError.Recipients[i] = recipient.Email
	}

	// Track successful sends in this batch
	var successCount int32 = 0

	// Process each recipient
	for _, recipient := range recipients {
		// Send the email with detailed error handling
		emailSent, sendErr := s.sendTransactionalEmailWithError(
			ctx,
			recipient.Email,
			getMessageContent(&domain.EmailRequest{
				HtmlContent: params.HtmlContent,
				Text:        params.Text,
			}),
			params.Subject,
			params.Sender.Email,
			getSenderName(params.Sender),
			params.UserID,
			params.CompanyID,
			&params.AuthUser,
		)

		if emailSent {
			successCount++
		} else {
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient.Email,
				Error: fmt.Sprintf("Failed to send email: %v", sendErr),
			})
		}
	}

	// Only update the database once per batch with the total count
	if successCount > 0 {
		mu.Lock()
		defer mu.Unlock()

		// Get current usage record
		mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, params.CompanyID)
		if err != nil {
			log.Printf("error fetching mail usage record: %v", err)
			// Add this as an error for all successfully sent emails in this batch
			for _, recipient := range recipients {
				found := false
				for _, existingErr := range batchError.Errors {
					if existingErr.Email == recipient.Email {
						found = true
						break
					}
				}
				if !found {
					batchError.Errors = append(batchError.Errors, RecipientError{
						Email: recipient.Email,
						Error: fmt.Sprintf("Email sent but failed to update usage record: %v", err),
					})
				}
			}
		} else {
			// Update with the count of successfully sent emails
			_, err = s.store.UpdateEmailsSentAndRemaining(ctx, db.UpdateEmailsSentAndRemainingParams{
				CompanyID:  params.CompanyID,
				EmailsSent: sql.NullInt32{Int32: successCount, Valid: true},
				ID:         mailUsageRecord.ID,
			})
			if err != nil {
				log.Printf("error updating mail usage: %v", err)
				// Add this as a warning for all successfully sent emails
				for _, recipient := range recipients {
					found := false
					for _, existingErr := range batchError.Errors {
						if existingErr.Email == recipient.Email {
							found = true
							break
						}
					}
					if !found {
						batchError.Errors = append(batchError.Errors, RecipientError{
							Email: recipient.Email,
							Error: fmt.Sprintf("Email sent but failed to update usage count: %v", err),
						})
					}
				}
			}
		}

		if len(batchError.Errors) == 0 {
			log.Printf("Successfully sent %d transactional emails in batch %d", successCount, batchIndex)
		} else {
			log.Printf("Transactional batch %d completed: %d sent, %d errors", batchIndex, successCount, len(batchError.Errors))
		}
	}

	return batchError
}

// sendTransactionalEmailWithError sends a single transactional email with detailed error reporting
func (s *MiscService) sendTransactionalEmailWithError(ctx context.Context, recipient string, emailContent string, subject string, from string, fromName string, userId uuid.UUID, companyId uuid.UUID, authUser *domain.SMTPAuthUser) (bool, error) {
	validEmail := helper.IsValidEmail(recipient)
	if !validEmail {
		return false, fmt.Errorf("invalid email address format")
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
		AuthUser:    *authUser,
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
			// The mail will be sent from the app's domain
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

func (s *MiscService) sendEmailWithSMTP(request *domain.EmailRequest) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	analyzer, err := helper.NewContentAnalyzer("internal/config/config.json", logger)
	if err != nil {
		return fmt.Errorf("failed to create content analyzer: %w", err)
	}

	// Get content for analysis
	content := getMessageContent(request)

	// Analyze the content before sending the email
	analysisResult, err := analyzer.AnalyzeContent(context.TODO(), content, nil)
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

func (s *MiscService) signEmail(domainEmail string, companyId uuid.UUID, emailBody []byte) ([]byte, error) {
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

func getMessageContent(emailRequest *domain.EmailRequest) string {
	if emailRequest.HtmlContent != nil {
		return *emailRequest.HtmlContent
	}
	if emailRequest.Text != nil {
		return *emailRequest.Text
	}
	return ""
}

func getSenderName(sender domain.Sender) string {
	if sender.Name != nil {
		return *sender.Name
	}
	return ""
}

func (s *MiscService) GetAllPlansWithDetails(ctx context.Context) ([]dto.PlanResponse, error) {
	// Use the optimized query that gets all plans with details
	plans, err := s.store.ListPlansWithDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting plans with details: %w", err)
	}

	var responses []dto.PlanResponse
	for _, plan := range plans {
		// Parse features from JSON or use separate query
		var planFeatures []dto.PlanFeature

		// For now, use separate query for features
		features, err := s.store.GetPlanFeaturesByPlanID(ctx, plan.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting plan features for plan %s: %w", plan.ID, err)
		}

		for _, feature := range features {
			planFeatures = append(planFeatures, dto.PlanFeature{
				Name:        feature.Name.String,
				Description: feature.Description.String,
				Value:       feature.Value.String,
			})
		}

		response := dto.PlanResponse{
			ID:           plan.ID,
			Name:         plan.Name,
			Description:  plan.Description.String,
			Price:        plan.Price.InexactFloat64(),
			BillingCycle: plan.BillingCycle.String,
			Status:       plan.Status.String,
			Features:     planFeatures,
			MailingLimits: dto.MailingLimits{
				DailyLimit:           plan.DailyLimit.Int32,
				MonthlyLimit:         plan.MonthlyLimit.Int32,
				MaxRecipientsPerMail: plan.MaxRecipientsPerMail.Int32,
			},
			CreatedAt: plan.CreatedAt.Format(time.RFC3339),
			UpdatedAt: plan.UpdatedAt.Format(time.RFC3339),
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *MiscService) GetSinglePlan(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error) {
	// Get plan by ID
	plan, err := s.store.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan by ID: %w", err)
	}

	// Get features
	features, err := s.store.GetPlanFeaturesByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan features: %w", err)
	}

	// Get mailing limits
	mailingLimit, err := s.store.GetMailingLimitByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting mailing limits: %w", err)
	}

	// Convert features
	var planFeatures []dto.PlanFeature
	for _, feature := range features {
		planFeatures = append(planFeatures, dto.PlanFeature{
			Name:        feature.Name.String,
			Description: feature.Description.String,
			Value:       feature.Value.String,
		})
	}

	// Prepare response
	response := &dto.PlanResponse{
		ID:           plan.ID,
		Name:         plan.Name,
		Description:  plan.Description.String,
		Price:        plan.Price.InexactFloat64(),
		BillingCycle: plan.BillingCycle.String,
		Status:       plan.Status.String,
		Features:     planFeatures,
		MailingLimits: dto.MailingLimits{
			DailyLimit:           mailingLimit.DailyLimit.Int32,
			MonthlyLimit:         mailingLimit.MonthlyLimit.Int32,
			MaxRecipientsPerMail: mailingLimit.MaxRecipientsPerMail.Int32,
		},
		CreatedAt: plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt: plan.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}
