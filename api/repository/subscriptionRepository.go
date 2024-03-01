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

func (r *SubscriptionRepository) GetAllSubscriptions() ([]model.SubscriptionResponseModel, error) {

	var subscriptions []model.SubscriptionResponseModel

	if err := r.DB.Preload("Plan").
		Preload("User").
		Preload("Billing").
		Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetAllCurrentRunningSubscription() ([]model.SubscriptionResponseModel, error) {

	var subscription model.Subscription

	if err := r.DB.Find(&subscription).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return []model.SubscriptionResponseModel{}, nil
		}

		return nil, fmt.Errorf("failed to get subscription records: %w", err)
	}

	return nil, nil
}

func (r *SubscriptionRepository) UpdateExpiredSubscription(id int) error {

	return nil
}

func (r *SubscriptionRepository) CancelSubscriptionService(id string, userId int) error {

	return nil
}

func (r *SubscriptionRepository) FindSubscriptionById(id string, userId int) (*model.SubscriptionResponseModel, error) {
	return nil, nil
}
