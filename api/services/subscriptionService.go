package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
)

type SubscriptionService struct {
	SubscriptionRepo *repository.SubscriptionRepository
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
