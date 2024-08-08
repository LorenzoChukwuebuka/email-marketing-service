package services

import (
	"email-marketing-service/api/v1/dto"
	smtpfactory "email-marketing-service/api/v1/factory/smtpFactory"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

type SMTPMailService struct {
	APIKeySVC        *APIKeyService
	SubscriptionRepo *repository.SubscriptionRepository
	DailyCalcRepo    *repository.DailyMailCalcRepository
	UserRepo         *repository.UserRepository
	MailStatusRepo   *repository.MailStatusRepository
	//mu               sync.Mutex
}

func NewSMTPMailService(apikeyservice *APIKeyService,
	subscriptionRepository *repository.SubscriptionRepository,
	dailyCalc *repository.DailyMailCalcRepository,
	userRepo *repository.UserRepository, mailStatusRepo *repository.MailStatusRepository) *SMTPMailService {
	return &SMTPMailService{
		APIKeySVC:        apikeyservice,
		SubscriptionRepo: subscriptionRepository,
		DailyCalcRepo:    dailyCalc,
		UserRepo:         userRepo,
		MailStatusRepo:   mailStatusRepo,
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

	// Get the daily mail calculator
	mailCalcRepo, err := s.DailyCalcRepo.GetDailyMailRecordForToday(int(subscription.ID))
	if err != nil {
		return nil, fmt.Errorf("error fetching record")
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
		if mailCalcRepo.RemainingMails == 0 {
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
			}, userId.ID)
		}()

		updateCalcData := &model.DailyMailCalc{
			UUID:           mailCalcRepo.UUID,
			RemainingMails: mailCalcRepo.RemainingMails - len(batch),
			MailsSent:      mailCalcRepo.MailsSent + len(batch),
		}

		err = s.DailyCalcRepo.UpdateDailyMailCalcRepository(updateCalcData)
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

	mailS, err := smtpfactory.MailFactory("mailtrap")
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
		num, err := strconv.Atoi(activeSub.Plan.NumberOfMailsPerDay)
		if err != nil {
			fmt.Println("Error converting NumberOfMailsPerDay to integer:", err)
			return err
		}
		
		dailyCalcData := &model.DailyMailCalc{
			UUID:           uuid.New().String(),
			SubscriptionID: int(activeSub.ID),
			MailsForADay:   num,
			MailsSent:      0,
			RemainingMails: num,
		}

		err = s.DailyCalcRepo.CreateRecordDailyMailCalculation(dailyCalcData)
		if err != nil {
			fmt.Println("Error creating daily mail calculation record:", err)
			return err
		}
	}

	return nil
}
