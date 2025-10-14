package worker

import (
	"context"
	"database/sql"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/helper"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unsafe"
	"github.com/google/uuid"
)

func (w *Worker) ProcessAPISMTPEmailsTask(ctx context.Context, payload SendAPISMTPEmailsPayload) error {

	mailUsageRecord, err := w.Store.GetCurrentEmailUsage(ctx, payload.CompanyId)
	if err != nil {
		return common.ErrFetchingRecord
	}

	if mailUsageRecord.RemainingEmails.Int32 == 0 {
		return fmt.Errorf("email quota exceeded")
	}

	authUser, err := w.Store.GetMasterSMTPKey(ctx, payload.UserId)
	if err != nil {
		log.Printf("error fetching master smtp key: %v", err)
		return fmt.Errorf("failed to fetch SMTP credentials: %w", err)
	}

	if authUser.Status != "active" {
		log.Printf("user status is not active")
		return fmt.Errorf("user status is not active")
	}

	authModel := &domain.SMTPAuthUser{
		Username: authUser.SmtpLogin,
		Password: authUser.Password,
	}

	// Check type of req.To and convert it to a slice of Recipient if necessary
	var recipients []domain.Recipient
	switch to := payload.EmailPayload.To.(type) {
	case domain.Recipient:
		recipients = []domain.Recipient{to}
	case []domain.Recipient:
		recipients = to
	default:
		return fmt.Errorf("invalid recipient type")
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

			batchErr := w.sendTransactionalEmailBatch(&mu, ctx, &TransactionalEmailParams{
				CompanyID:   payload.CompanyId,
				UserID:      payload.UserId,
				AuthUser:    *authModel,
				Sender:      payload.EmailPayload.Sender,
				Subject:     payload.EmailPayload.Subject,
				HtmlContent: payload.EmailPayload.HtmlContent,
				Text:        payload.EmailPayload.Text,
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

	_ = result

	return nil

}

type TransactionalEmailParams struct {
	CompanyID   uuid.UUID
	UserID      uuid.UUID
	AuthUser    domain.SMTPAuthUser
	Sender      domain.Sender
	Subject     string
	HtmlContent *string
	Text        *string
}

func (s *Worker) sendTransactionalEmailBatch(mu *sync.Mutex, ctx context.Context, params *TransactionalEmailParams, recipients []domain.Recipient, batchIndex int) BatchError {
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
		mailUsageRecord, err := s.Store.GetCurrentEmailUsage(ctx, params.CompanyID)
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
			_, err = s.Store.UpdateEmailsSentAndRemaining(ctx, db.UpdateEmailsSentAndRemainingParams{
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

// func (s *Worker) signEmail_(domainEmail string, companyId uuid.UUID, emailBody []byte) ([]byte, error) {
// 	// Fetch the domain associated with the sender
// 	domain, err := s.Store.FindDomainByNameAndCompany(context.Background(), db.FindDomainByNameAndCompanyParams{
// 		Domain:    domainEmail,
// 		CompanyID: companyId,
// 	})

// 	if err != nil || !domain.Verified.Valid {
// 		return nil, fmt.Errorf("domain not found or not verified")
// 	}

// 	helper.ValidatePrivateKey(domain.DkimPrivateKey.String)

// 	// DKIM signing process
// 	signedEmail, err := helper.SignEmail(&emailBody, domain.Domain, domain.DkimSelector.String, string(domain.DkimPrivateKey.String))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to sign email: %v", err)
// 	}

// 	return signedEmail, nil
// }

// sendTransactionalEmailWithError sends a single transactional email with detailed error reporting
func (s *Worker) sendTransactionalEmailWithError(ctx context.Context, recipient string, emailContent string, subject string, from string, fromName string, userId uuid.UUID, companyId uuid.UUID, authUser *domain.SMTPAuthUser) (bool, error) {
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

	sender_domain, err := s.Store.FindDomainByNameAndCompany(ctx, db.FindDomainByNameAndCompanyParams{
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
