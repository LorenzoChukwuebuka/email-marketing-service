package services

import (
	"email-marketing-service/api/dto"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type SMTPMailService struct {
	APIKeySVC        *APIKeyService
	SubscriptionRepo *repository.SubscriptionRepository
	DailyCalcRepo    *repository.DailyMailCalcRepository
	UserRepo         *repository.UserRepository
}

func NewSMTPMailService(apikeyservice *APIKeyService,
	subscriptionRepository *repository.SubscriptionRepository,
	dailyCalc *repository.DailyMailCalcRepository,
	userRepo *repository.UserRepository) *SMTPMailService {
	return &SMTPMailService{
		APIKeySVC:        apikeyservice,
		SubscriptionRepo: subscriptionRepository,
		DailyCalcRepo:    dailyCalc,
		UserRepo:         userRepo,
	}
}

const batchSize = 20 // Adjust this value as per your requirements

func (s *SMTPMailService) SendSMTPMail(d *dto.EmailRequest, apiKey string) (map[string]interface{}, error) {
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
	mailCalcRepo, err := s.DailyCalcRepo.GetDailyMailRecordForToday(subscription.Id)
	if err != nil {
		return nil, fmt.Errorf("error fetching record")
	}

	// Calculate the number of batches
	numBatches := len(d.To) / batchSize

	if len(d.To)%batchSize != 0 {
		numBatches++
	}

	// Process emails in batches
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(d.To) {
			end = len(d.To)
		}
		batch := d.To[start:end]

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
			})
		}()

		updateCalcData := &model.DailyMailCalc{
			UUID:           mailCalcRepo.UUID,
			RemainingMails: mailCalcRepo.RemainingMails - len(batch),
			MailsSent:      mailCalcRepo.MailsSent + len(batch),
		}

		// Update the counter
		mailCalcRepo.RemainingMails -= len(batch)
		mailCalcRepo.MailsSent += len(batch)
		err = s.DailyCalcRepo.UpdateDailyMailCalcRepository(updateCalcData)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *SMTPMailService) handleSendMail(emailRequest *dto.EmailRequest) {
	for _, recipient := range emailRequest.To {
		time.Sleep(2 * time.Second)
		fmt.Printf("Mail sent to %s\n", recipient.Email)
	}
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

		println(num)

		dailyCalcData := &model.DailyMailCalc{
			UUID:           uuid.New().String(),
			SubscriptionID: activeSub.Id,
			MailsForADay:   num,
			MailsSent:      0,
			CreatedAt:      time.Now(),
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
