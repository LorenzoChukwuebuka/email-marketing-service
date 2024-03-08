package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
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
	return nil, nil
}
