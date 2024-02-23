package repository

import (

	"email-marketing-service/api/model"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) CreateSubscription(d *model.Subscription) error {
	 return nil
}

func (r *SubscriptionRepository) GetAllSubscriptions() ([]model.SubscriptionResponseModel, error) {
	 

	return nil, nil
}

func (r *SubscriptionRepository) GetAllCurrentRunningSubscription() ([]model.SubscriptionResponseModel, error) {
	

	return nil, nil
}

func (r *SubscriptionRepository) UpdateExpiredSubscription(id int) error {
	 

	return nil
}

func (r *SubscriptionRepository) CancelSubscriptionService(id string, userId int) error {
	 

	return nil
}

func (r *SubscriptionRepository) FindSubscriptionById(id string, userId int) (*model.SubscriptionResponseModel, error) {
	 return nil,nil
}
