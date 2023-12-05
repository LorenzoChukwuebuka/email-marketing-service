package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"fmt"
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

/*
These are mostly jobs
*/

func (s *SubscriptionService) GetAllSubscription() ([]model.SubscriptionResponseModel, error) {
	subscriptions, err := s.SubscriptionRepo.GetAllSubscriptions()

	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	var expiredUserIDs []int

	for _, subscription := range subscriptions {
		if subscription.EndDate.Before(currentTime) {
			fmt.Printf("Subscription ID %d has expired.\n", subscription.Id)
			expiredUserIDs = append(expiredUserIDs, subscription.UserId)
		} else {
			fmt.Printf("Subscription ID %d has not expired.\n", subscription.Id)

		}

	}

	if len(expiredUserIDs) != 0 {

		for _, userIds := range expiredUserIDs {
			fmt.Println(userIds)
		}

	}

	fmt.Println(expiredUserIDs)

	return subscriptions, err
}

func (s *SubscriptionService) SendSubscriptionExpiryNotificationReminder() {}
