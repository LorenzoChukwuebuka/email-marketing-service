package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) createSubscriptionResponse(s model.Subscription) *model.SubscriptionResponseModel {

	response := &model.SubscriptionResponseModel{
		UUID:      s.UUID,
		UserId:    s.UserId,
		PlanId:    s.PlanId,
		PaymentId: s.PaymentId,
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
		Expired:   s.Expired,
		CreatedAt: s.CreatedAt,
	}

	if s.UpdatedAt != nil {
		response.UpdatedAt = s.UpdatedAt.Format(time.RFC3339)
	} else {
		response.UpdatedAt = ""
	}

	if s.Plan != nil {
		response.Plan = &model.PlanResponse{
			UUID:                s.Plan.UUID,
			PlanName:            s.Plan.PlanName,
			Duration:            s.Plan.Duration,
			Price:               s.Plan.Price,
			NumberOfMailsPerDay: s.Plan.NumberOfMailsPerDay,
			Details:             s.Plan.Details,
			Status:              s.Plan.Status,
			CreatedAt:           s.Plan.CreatedAt,
		}

		if s.Plan.UpdatedAt != nil {
			response.Plan.UpdatedAt = s.Plan.UpdatedAt.Format(time.RFC3339)
		} else {
			response.Plan.UpdatedAt = ""
		}

	}

	if s.User != nil {
		response.User = &model.UserResponse{
			UUID:     s.User.UUID,
			FullName: s.User.FullName,
		}
	}

	return response

}

func (r *SubscriptionRepository) CreateSubscription(d *model.Subscription) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) GetAllSubscriptions() ([]model.Subscription, error) {

	var subscriptions []model.Subscription

	if err := r.DB.Preload("Plan").
		Preload("User").
		Preload("Billing").
		Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetAllCurrentRunningSubscription() ([]model.Subscription, error) {

	var subscription []model.Subscription

	if err := r.DB.Where("expired = ?", false).Find(&subscription).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return []model.Subscription{}, nil
		}

		return nil, fmt.Errorf("failed to get subscription records: %w", err)
	}

	return subscription, nil
}

func (r *SubscriptionRepository) UpdateExpiredSubscription(id int) error {
	// Fetch the subscription by ID
	var subscription model.Subscription
	if err := r.DB.First(&subscription, id).Error; err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	// Update the 'expired' field to true
	subscription.Expired = true

	// Save the updated subscription
	if err := r.DB.Save(&subscription).Error; err != nil {
		return fmt.Errorf("failed to update expired subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) CancelSubscriptionService(id string, userId int) error {
	var subscription model.Subscription
	if err := r.DB.Where("id = ? AND user_id = ?", id, userId).First(&subscription).Error; err != nil {
		return fmt.Errorf("failed to find subscription for cancellation: %w", err)
	}

	subscription.Cancelled = true

	// Save the updated subscription (e.g., update cancellation status)
	if err := r.DB.Save(&subscription).Error; err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) FindSubscriptionById(id string, userId int) (*model.SubscriptionResponseModel, error) {
	var subscription model.Subscription
	if err := r.DB.Preload("Plan").Preload("User").Preload("Billing").Where("id = ? AND user_id = ?", id, userId).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}
	return r.createSubscriptionResponse(subscription), nil
}

//for the future....

// func (r *SubscriptionRepository) GetAllSubscriptions() ([]model.Subscription, error) {
// 	var subscriptions []model.Subscription

// 	if err := r.DB.Preload("Plan").
// 		Preload("User").
// 		Preload("Billing", func(db *gorm.DB) *gorm.DB {
// 			return db.Preload("Plan").
// 				Preload("User")
// 		}).
// 		Find(&subscriptions).Error; err != nil {
// 		return nil, fmt.Errorf("failed to fetch subscriptions: %w", err)
// 	}

// 	return subscriptions, nil
// }
