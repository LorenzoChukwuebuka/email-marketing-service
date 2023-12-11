package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
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

	fmt.Println(userId)

	return nil, nil
}

func (s *SMTPMailService) CreateRecordForDailyMailCalculation() error {
	//1.... Select all active subscriptions....

	activeSubs, err := s.SubscriptionRepo.GetAllCurrentRunningSubscription()
	if err != nil {
		return err
	}

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
		}

		err = s.DailyCalcRepo.CreateRecordDailyMailCalculation(dailyCalcData)

		if err != nil {
			return err
		}
	}

	return nil

}
