package services

import (
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
}

func NewSMTPMailService(apikeyservice *APIKeyService, subscriptionRepository *repository.SubscriptionRepository, dailyCalc *repository.DailyMailCalcRepository) *SMTPMailService {
	return &SMTPMailService{
		APIKeySVC:        apikeyservice,
		SubscriptionRepo: subscriptionRepository,
		DailyCalcRepo:    dailyCalc,
	}
}

func (s *SMTPMailService) SendSMTPMail(d *model.EmailRequest, apiKey string) (map[string]interface{}, error) {

	//1. Get user's Id
	userId, err := s.APIKeySVC.FindUserWithAPIKey(apiKey)

	if err != nil {
		return nil, fmt.Errorf("error fetching userId")
	}

	//2. Get the daily mail calculator

	mailCalcRepo, err := s.DailyCalcRepo.GetDailyMailRecordForToday(userId)

	if err != nil {
		return nil, fmt.Errorf("error fetching sent mails")
	}

	//3. check remaining mails if it is equals to 0

	if mailCalcRepo.RemainingMails == 0 {
		return nil, fmt.Errorf("you have exceeded your daily plan")
	}

	mailResult := make(chan interface{})
	go s.handleSendMail(mailResult)

	//4. Handle the result from handleSendMail if needed
	result := <-mailResult
	fmt.Printf("Mail result: %+v\n", result)

	//5. update counter

	updateCalcData := &model.DailyMailCalcModel{
		UUID:           mailCalcRepo.UUID,
		RemainingMails: mailCalcRepo.RemainingMails - 1,
		MailsSent:      mailCalcRepo.MailsSent + 1,
	}

	err = s.DailyCalcRepo.UpdateDailyMailCalcRepository(updateCalcData)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *SMTPMailService) handleSendMail(resultChan chan interface{}) {
	// Perform the mail sending logic here

	// For example, simulate sending mail for demonstration purposes
	time.Sleep(2 * time.Second)

	// Send the result to the channel
	resultChan <- "Mail sent successfully"
}

//##################################################### JOBS #################################################################

func (s *SMTPMailService) CreateRecordForDailyMailCalculation() error {
	//1.... Select all active subscriptions....

	activeSubs, err := s.SubscriptionRepo.GetAllCurrentRunningSubscription()
	if err != nil {
		return err
	}

	fmt.Println(activeSubs)

	for _, activeSub := range activeSubs {

		num, err := strconv.Atoi(activeSub.Plan.NumberOfMailsPerDay)

		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		dailyCalcData := &model.DailyMailCalcModel{
			UUID:           uuid.New().String(),
			SubscriptionID: activeSub.Id,
			MailsForADay:   num,
			MailsSent:      0,
			CreatedAt:      time.Now(),
			RemainingMails: num,
		}

		err = s.DailyCalcRepo.CreateRecordDailyMailCalculation(dailyCalcData)

		if err != nil {
			return err
		}

	}

	return nil

}
