package repository

import "database/sql"

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) CreateSubscription() error {
	return nil
}

func (r *SubscriptionRepository) GetAllSubscription() {}


func (r *SubscriptionRepository) GetExpiredSubscriptions(){}

func (r *SubscriptionRepository) GetSingleSubscription(id int){}
