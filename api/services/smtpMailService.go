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

	// Get the daily mail calculator

	mailCalcRepo, err := s.DailyCalcRepo.GetDailyMailRecordForToday(userId.ID)

	if err != nil {
		return nil, fmt.Errorf("error fetching record")
	}

	//. check remaining mails if it is equals to 0

	if mailCalcRepo.RemainingMails == 0 {
		return nil, fmt.Errorf("you have exceeded your daily plan")
	}

	mailResult := make(chan interface{}, 1)
	go s.handleSendMail(mailResult)

	//. Handle the result from handleSendMail if needed
	result := <-mailResult
	fmt.Printf("Mail result: %+v\n", result)

	//. update counter

	updateCalcData := &model.DailyMailCalc{
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
			return err
		}

	}

	return nil

}
