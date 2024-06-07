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

func (s *SubscriptionService) CreateSubscription(d *model.Subscription) (*model.Subscription, error) {
	//check if the user is running on a free plan and expire it automatically

	if err := s.SubscriptionRepo.CreateSubscription(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *SubscriptionService) CancelSubscriptionService(userId int, subscriptionId string) (map[string]interface{}, error) {
	//0.5 check if the user already cancelled
	userCancelledSub, err := s.SubscriptionRepo.FindSubscriptionById(subscriptionId, userId)
	if err != nil {
		return nil, err
	}

	if userCancelledSub.Cancelled {
		return nil, fmt.Errorf("subscription already cancelled")
	}

	// 1. The user cancels the subscription
	err = s.SubscriptionRepo.CancelSubscriptionService(subscriptionId, userId)
	if err != nil {
		return nil, err
	}

	// 2. Get the current running subscription of the user
	userSub, err := s.SubscriptionRepo.FindSubscriptionById(subscriptionId, userId)
	if err != nil {
		return nil, err
	}

	// 3. Calculate how much duration is left of their subscription
	layout := "2006-01-02T15:04:05.999999-07:00"
	timeStr1 := userSub.DateCancelled
	timeStr2 := userSub.EndDate.Format(layout)

	t1, err := time.Parse(layout, timeStr1)
	if err != nil {
		return nil, fmt.Errorf("error parsing timeStr1: %w", err)
	}

	t2, err := time.Parse(layout, timeStr2)
	if err != nil {
		return nil, fmt.Errorf("error parsing timeStr2: %w", err)
	}

	// 4. Calculate the duration between the two times
	duration := t2.Sub(t1)
	remainingDays := int(duration.Hours() / 24)

	//5. Handles the calculation for the amount to refund
	amountToRefund, err := calculateAmountToRefund(remainingDays, userSub.StartDate, userSub.EndDate, userSub.Billing.AmountPaid)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	fmt.Println(amountToRefund)

	//initiate the refund policy from here but for now I will just have to output it to the user

	successMap := map[string]interface{}{
		"status":         "refunded successfully",
		"amountToRefund": amountToRefund,
	}

	return successMap, nil
}

func calculateAmountToRefund(remainingDays int, startDate time.Time, endDate time.Time, amountPaid float32) (float32, error) {

	//1. total number of days
	totalDuration := endDate.Sub(startDate).Hours() / 24
	if totalDuration <= 0 {
		return 0, fmt.Errorf("invalid duration between start and end date")
	}

	//2. Amount per day

	amountPerDay := amountPaid / float32(totalDuration)

	//3. Amount to refund

	amountToRefund := amountPerDay * float32(remainingDays)

	return float32(amountToRefund), nil
}

/** ############################################################### JOBS ####################################################################### **/
func (s *SubscriptionService) UpdateExpiredSubscription() ([]model.Subscription, error) {
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
			err = mail.SubscriptionExpiryMail(subscription.User.FullName, subscription.User.Email, subscription.Plan.PlanName)
			if err != nil {
				return nil, err
			}
		} else {
			fmt.Printf("Subscription ID %d has not expired.\n", subscription.Id)
		}
	}

	return subscriptions, err
}

func (s *SubscriptionService) SendSubscriptionExpiryNotificationReminder() error {
	subscriptions, err := s.SubscriptionRepo.GetAllCurrentRunningSubscription()

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
