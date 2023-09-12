package repository

import "database/sql"

type SubscriptionRepository struct {
	DB *sql.DB
}


 func (r *SubscriptionRepository) CreateSubscription() {}