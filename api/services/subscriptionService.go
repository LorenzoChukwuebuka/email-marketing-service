package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"time"
)

type SubscriptionService struct {
	SubscriptionRepo *repository.SubscriptionRepository
}

type CurrentSubscription struct {
	ID        int
	UUID      string
	UserID    int
	PlanID    int
	PaymentID int
	StartDate time.Time
	EndDate   time.Time
	Expired   bool
}

func NewSubscriptionService(subscriptionRepo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{SubscriptionRepo: subscriptionRepo}
}

func (s *SubscriptionService) CreateSubscription(d *model.SubscriptionModel) (*model.SubscriptionModel, error) {

	if err := s.SubscriptionRepo.CreateSubscription(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *SubscriptionService) GetAllCurrentRunningSubscription(subscriptionId string) ([]model.SubscriptionModel, error) {
	subscriptions, err := s.SubscriptionRepo.GetAllSubscriptions(subscriptionId)

	if err != nil {
		return nil, err
	}

	return subscriptions, err
}

/*
These are mostly jobs

*/

func (s *SubscriptionService) SendSubscriptionExpiryNotificationReminder() {}
