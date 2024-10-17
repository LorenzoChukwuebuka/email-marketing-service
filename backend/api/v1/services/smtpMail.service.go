package services

import (
	"email-marketing-service/api/v1/dto"
	smtpfactory "email-marketing-service/api/v1/factory/smtpFactory"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type SMTPMailService struct {
	APIKeySVC        *APIKeyService
	SubscriptionRepo *repository.SubscriptionRepository
	MailUsageRepo    *repository.MailUsageRepository
	UserRepo         *repository.UserRepository
	MailStatusRepo   *repository.MailStatusRepository
	PlanRepo         *repository.PlanRepository
	SMTPKeyRepo      *repository.SMTPKeyRepository
	//mu               sync.Mutex
}

func NewSMTPMailService(apikeyservice *APIKeyService,
	subscriptionRepository *repository.SubscriptionRepository,
	mailUsageRepo *repository.MailUsageRepository,
	userRepo *repository.UserRepository, mailStatusRepo *repository.MailStatusRepository, planRepo *repository.PlanRepository, smtpkeyRepo *repository.SMTPKeyRepository) *SMTPMailService {
	return &SMTPMailService{
		APIKeySVC:        apikeyservice,
		SubscriptionRepo: subscriptionRepository,
		MailUsageRepo:    mailUsageRepo,
		UserRepo:         userRepo,
		MailStatusRepo:   mailStatusRepo,
		PlanRepo:         planRepo,
		SMTPKeyRepo:      smtpkeyRepo,
	}
}

const batchSize = 20

func (s *SMTPMailService) PrepareMail(d *dto.EmailRequest, apiKey string) (map[string]interface{}, error) {
	userUUID, err := s.APIKeySVC.FindUserWithAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("error fetching userId")
	}
	userModel := &model.User{UUID: userUUID}
	userId, err := s.UserRepo.FindUserById(userModel)
	if err != nil {
		return nil, fmt.Errorf("error fetching userId")
	}

	// Get the current user's running subscription
	subscription, err := s.SubscriptionRepo.GetUserCurrentRunningSubscription(userId.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching subscription record")
	}

	// Get the user's mail usage for today
	mailUsageRepo, err := s.MailUsageRepo.GetCurrentMailUsageRecord(int(subscription.ID))
	if err != nil {
		return nil, fmt.Errorf("error fetching mail usage record")
	}

	//get the master key from the db

	authUser, err := s.SMTPKeyRepo.GetSMTPMasterKey(userUUID)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	authModel := &dto.SMTPAuthUser{
		Username: authUser.SMTPLogin,
		Password: authUser.Password,
	}

	// Check type of d.To and convert it to a slice of Recipient if necessary
	var recipients []dto.Recipient
	switch to := d.To.(type) {
	case dto.Recipient:
		recipients = []dto.Recipient{to}
	case []dto.Recipient:
		recipients = to
	default:
		return nil, fmt.Errorf("invalid recipient type")
	}

	// Calculate the number of batches
	numBatches := len(recipients) / batchSize
	if len(recipients)%batchSize != 0 {
		numBatches++
	}

	// Process emails in batches
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(recipients) {
			end = len(recipients)
		}
		batch := recipients[start:end]

		// Check remaining mails
		if mailUsageRepo.RemainingMails == 0 {
			return nil, fmt.Errorf("you have exceeded your daily plan")
		}

		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Error sending batch: %v\n", err)
				}
			}()
			go s.handleSendMail(&dto.EmailRequest{
				Sender:      d.Sender,
				To:          batch,
				Subject:     d.Subject,
				HtmlContent: d.HtmlContent,
				Text:        d.Text,
				AuthUser:    *authModel,
			}, userId.ID)
		}()

		updateUsageData := &model.MailUsage{
			UUID:           mailUsageRepo.UUID,
			RemainingMails: mailUsageRepo.RemainingMails - len(batch),
			MailsSent:      mailUsageRepo.MailsSent + len(batch),
		}

		err = s.MailUsageRepo.UpdateMailUsageRecord(updateUsageData)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *SMTPMailService) handleSendMail(emailRequest *dto.EmailRequest, userId uint) error {
	var recipients []dto.Recipient
	switch to := emailRequest.To.(type) {
	case dto.Recipient:
		recipients = []dto.Recipient{to}
	case []dto.Recipient:
		recipients = to
	default:
		return fmt.Errorf("invalid recipient type")
	}

	for _, recipient := range recipients {
		UUID := uuid.New().String()

		emailData := &model.SentEmails{
			UUID:           UUID,
			Sender:         uint(userId),
			Recipient:      recipient.Email,
			MessageContent: getMessageContent(emailRequest),
			Status:         model.Sending,
		}

		err := s.MailStatusRepo.CreateStatus(emailData)
		if err != nil {
			return err
		}
	}

	mailS, err := smtpfactory.MailFactory(config.MAIL_PROCESSOR)
	if err != nil {
		// Update the status to "failed" for the existing records
		for _, recipient := range recipients {
			emailData, err := s.MailStatusRepo.GetSentEmailByRecipient(recipient.Email)
			if err != nil {
				return err
			}
			emailData.Status = model.Failed
			s.MailStatusRepo.UpdateReport(emailData)
		}
		return err
	}

	err = mailS.HandleSendMail(emailRequest)
	if err != nil {
		// Update the status to "failed" for the existing records
		for _, recipient := range recipients {
			emailData, err := s.MailStatusRepo.GetSentEmailByRecipient(recipient.Email)
			if err != nil {
				return err
			}
			emailData.Status = model.Failed
			s.MailStatusRepo.UpdateReport(emailData)
		}
		return err
	}

	// Update the status to "success" for the existing records
	for _, recipient := range recipients {
		emailData, err := s.MailStatusRepo.GetSentEmailByRecipient(recipient.Email)
		if err != nil {
			return err
		}
		emailData.Status = model.Success
		s.MailStatusRepo.UpdateReport(emailData)
	}

	return nil
}

func getMessageContent(emailRequest *dto.EmailRequest) string {
	if emailRequest.HtmlContent != nil {
		return *emailRequest.HtmlContent
	}
	if emailRequest.Text != nil {
		return *emailRequest.Text
	}
	return ""
}

//##################################################### JOBS #################################################################

func (s *SMTPMailService) CreateRecordForDailyMailCalculation() error {
	// Select all active subscriptions
	activeSubs, err := s.SubscriptionRepo.GetAllCurrentRunningSubscription()
	if err != nil {
		fmt.Println("Error fetching subscriptions:", err)
		return err
	}

	if len(activeSubs) == 0 {
		fmt.Println("No active subscriptions found.")
		return nil
	}

	for _, activeSub := range activeSubs {
		// Fetch the plan's mailing limit
		plan, err := s.PlanRepo.GetPlanById(activeSub.PlanId)
		if err != nil {
			fmt.Printf("Error fetching plan for subscription %d: %v\n", activeSub.ID, err)
			continue
		}

		// Only create records for plans with daily limits
		if plan.MailingLimit.LimitPeriod != "day" {
			continue
		}

		// Create a new mail usage record
		now := time.Now()
		periodStart := now.Truncate(24 * time.Hour)                    // Start of the day (12:00 AM)
		periodEnd := periodStart.Add(24 * time.Hour).Add(-time.Second) // End of the day (11:59:59 PM)

		mailUsageData := &model.MailUsage{
			UUID:           uuid.New().String(),
			SubscriptionID: int(activeSub.ID),
			PeriodStart:    periodStart,
			PeriodEnd:      periodEnd,
			LimitAmount:    plan.MailingLimit.LimitAmount,
			MailsSent:      0,
		}

		err = s.MailUsageRepo.CreateMailUsageRecord(mailUsageData)
		if err != nil {
			fmt.Printf("Error creating mail usage record for subscription %d: %v\n", activeSub.ID, err)
			continue
		}
	}

	return nil
}
