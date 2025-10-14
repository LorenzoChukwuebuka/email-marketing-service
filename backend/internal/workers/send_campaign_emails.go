package worker

import (
	"bytes"
	"context"
	"database/sql"
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

const BATCH_SIZE = 100

var (
	cfg = config.LoadEnv()
)

func (p *Worker) ProcessSendCampaignEmailsTask(ctx context.Context, payload SendCampaignEmailPayload) error {
	//get all associated emails for this campaign (from the contact groups)
	contactEmails, err := p.Store.GetCampaignContactEmails(ctx, payload.CampaignID)
	if err != nil {
		log.Printf("error occurred while fetching contacts: %v", err)
		_ = CampaignError{
			Type:    "CONTACT_FETCH_ERROR",
			Message: "Failed to fetch campaign contacts",
			Err:     err,
		}
	}

	if len(contactEmails) == 0 {
		_ = CampaignError{
			Type:    "NO_CONTACTS_ERROR",
			Message: "No contacts found for campaign",
			Err:     nil,
		}
	}

	mailUsageRecord, err := p.Store.GetCurrentEmailUsage(ctx, payload.CompanyID)
	if err != nil {
		_ = CampaignError{
			Type:    "USAGE_CHECK_ERROR",
			Message: "Failed to fetch email usage record",
			Err:     err,
		}
	}

	if mailUsageRecord.RemainingEmails.Int32 == 0 {
		_ = CampaignError{
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

			batchErr := p.sendEmailBatch(&mu, ctx, &payload, batch, batchIdx)
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
	if err = p.Store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
		Status: sql.NullString{String: string(enums.CampaignStatus(enums.Sent)), Valid: true},
		ID:     payload.CampaignID,
		UserID: payload.UserID,
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

func (p *Worker) sendEmailBatch(mu *sync.Mutex, ctx context.Context, d *SendCampaignEmailPayload, recipients []string, batchIndex int) BatchError {
	batchError := BatchError{
		BatchIndex: batchIndex,
		Recipients: recipients,
		Errors:     []RecipientError{},
	}

	//first thing fetch the campaign and extract the template id
	campaign, err := p.Store.GetCampaignByID(ctx, db.GetCampaignByIDParams{
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

	user, err := p.Store.GetUserByID(ctx, d.UserID)
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
		modifiedTemplate, err := p.addTrackingToTemplate(campaign.TemplateEmailHtml.String, campaign.ID.String(), recipient, user.CompanyID.String())
		if err != nil {
			log.Printf("error adding tracking to template for recipient %s: %v", recipient, err)
			batchError.Errors = append(batchError.Errors, RecipientError{
				Email: recipient,
				Error: fmt.Sprintf("Failed to add tracking to template: %v", err),
			})
			continue
		}

		// Send the email and track success with detailed error
		emailSent, sendErr := p.sendEmailWithError(
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
		if err := p.createEmailCampaignResult(campaign.ID, d.CompanyID, recipient); err != nil {
			// This is not a critical error, but we should log it
			log.Printf("Warning: Failed to create campaign result for %s: %v", recipient, err)
		}
	}

	// Only update the database once per batch with the total count
	if successCount > 0 {
		mu.Lock()
		defer mu.Unlock() // Ensure unlock happens even if there's an error

		// Get current usage record
		mailUsageRecord, err := p.Store.GetCurrentEmailUsage(ctx, d.CompanyID)
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
			_, err = p.Store.UpdateEmailsSentAndRemaining(ctx, db.UpdateEmailsSentAndRemainingParams{
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

func (s *Worker) createEmailCampaignResult(campaignId uuid.UUID, companyId uuid.UUID, recipientEmail string) error {
	_, err := s.Store.CreateEmailCampaignResult(context.Background(), db.CreateEmailCampaignResultParams{
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

func (s *Worker) addTrackingToTemplate(template string, campaignId string, recipientEmail string, companyId string) (string, error) {
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
func (s *Worker) sendEmailWithError(ctx context.Context, recipient string, emailContent string, subject string, previewText string, from string, fromName string, userId uuid.UUID, companyId uuid.UUID) (bool, error) {
	validEmail := helper.IsValidEmail(recipient)

	if !validEmail {
		return false, fmt.Errorf("invalid email address format")
	}

	authUser, err := s.Store.GetMasterSMTPKey(ctx, userId)
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

	sender_domain, err := s.Store.FindDomainByNameAndCompany(ctx, db.FindDomainByNameAndCompanyParams{
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

func (s *Worker) sendEmailWithSMTP(request *domain.EmailRequest) error {
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

func (s *Worker) signEmail(domainEmail string, companyId uuid.UUID, emailBody []byte) ([]byte, error) {
	// Fetch the domain associated with the sender
	domain, err := s.Store.FindDomainByNameAndCompany(context.Background(), db.FindDomainByNameAndCompanyParams{
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

func (s *Worker) StoreCampaignError(ctx context.Context, campaignID uuid.UUID, campaignErr CampaignError) {
	errMessage := campaignErr.Message

	if campaignErr.Err != nil {
		errMessage = fmt.Sprintf("%s: %v", campaignErr.Message, campaignErr.Err)
	}

	dbErr := s.Store.CreateCampaignError(ctx, db.CreateCampaignErrorParams{
		CampaignID:   campaignID,
		ErrorType:    sql.NullString{String: campaignErr.Type, Valid: true},
		ErrorMessage: errMessage,
	})
	if dbErr != nil {
		log.Printf("Failed to store campaign error: %v", dbErr)
	}
}
