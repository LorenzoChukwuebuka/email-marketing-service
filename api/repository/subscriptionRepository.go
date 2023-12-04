package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) CreateSubscription(d *model.SubscriptionModel) error {
	query := "INSERT INTO subscriptions (uuid,user_id, plan_id, payment_id, start_date, end_date, expired, created_at,transaction_id) VALUES ($1, $2, $3, $4, $5, $6, $7,$8,$9) RETURNING id"

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var insertedID int

	err = stmt.QueryRow(
		d.UUID,
		d.UserId,
		d.PlanId,
		d.PaymentId,
		d.StartDate,
		d.EndDate,
		d.Expired,
		d.CreatedAt,
		d.TransactionId,
	).Scan(&insertedID)

	if err != nil {
		return err
	}

	d.Id = insertedID
	return nil // Return nil here on success
}

func (r *SubscriptionRepository) GetAllSubscriptions(subscriptionId string) ([]model.SubscriptionModel, error) {
	query := `
      SELECT *
      FROM subscriptions;
    `

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subscriptions []model.SubscriptionModel

	for rows.Next() {
		var subscription model.SubscriptionModel

		err := rows.Scan(
			&subscription.Id,
			&subscription.UUID,
			&subscription.UserId,
			&subscription.PlanId,
			&subscription.PaymentId,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.Expired,
			&subscription.TransactionId,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, err
}

func (r *SubscriptionRepository) GetAllCurrentRunningSubscription() ([]model.SubscriptionModel, error) {
	query := `
      SELECT *
      FROM subscriptions
	  WHERE 
	  expired = FALSE
    `

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subscriptions []model.SubscriptionModel

	for rows.Next() {
		var subscription model.SubscriptionModel

		err := rows.Scan(
			&subscription.Id,
			&subscription.UUID,
			&subscription.UserId,
			&subscription.PlanId,
			&subscription.PaymentId,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.Expired,
			&subscription.TransactionId,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, err
}

func (r *SubscriptionRepository) GetCurrentSubscription(id int) {}
