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
	//check if the user is running on a free plan and expire it automatically

	if err := s.SubscriptionRepo.CreateSubscription(d); err != nil {
		return nil, err
	}
	return d, nil
}

/*
These are mostly jobs
*/

func (s *SubscriptionService) UpdateExpiredSubscription() ([]model.SubscriptionResponseModel, error) {
	subscriptions, err := s.SubscriptionRepo.GetAllSubscriptions()

	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	for _, subscription := range subscriptions {
		if subscription.EndDate.Before(currentTime) && !subscription.Expired {
			fmt.Printf("Subscription ID %d has expired.\n", subscription.Id)

			// Update the subscription as expired
			err := s.SubscriptionRepo.UpdateExpiredSubscription(subscription.Id)
			if err != nil {
				return nil, err
			}

			// Send mail to the user notifying them of the expiration of the service
			err = mail.SubscriptionExpiryMail(subscription.User.UserName, subscription.User.Email, subscription.Plan.PlanName)
			if err != nil {
				return nil, err
			}
		} else {
			fmt.Printf("Subscription ID %d has not expired.\n", subscription.Id)
		}
	}

	return subscriptions, err
}

func (s *SubscriptionService) CancelSubscriptionService(userId int) error {
	/**
	1. The user cancles the subscription
	2. A calculation is done which calculates how much is left of their subscription
	3. A refund is made after 24 hours automatically

	**/
	return nil
}

func (s *SubscriptionService) SendSubscriptionExpiryNotificationReminder() error {
	subscriptions, err := s.SubscriptionRepo.GetAllSubscriptions()

	if err != nil {
		return err
	}

	currentTime := time.Now()

	var expiredUserIDs []int

	for _, subscription := range subscriptions {
		// Calculate the time difference between current time and EndDate
		timeDifference := subscription.EndDate.Sub(currentTime)

		// Check if the time difference is 5 days or fewer
		if timeDifference.Hours() <= 24*5 {
			fmt.Printf("Subscription ID %d will expire within 5 days.\n", subscription.Id)
			expiredUserIDs = append(expiredUserIDs, subscription.UserId)
		} else {
			fmt.Printf("Subscription ID %d has not expired.\n", subscription.Id)
		}
	}

	fmt.Println(expiredUserIDs)

	return nil

}
