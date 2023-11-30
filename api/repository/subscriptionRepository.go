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

func (r *SubscriptionRepository) CheckExpiredSubscriptions(subscriptionId string) {
	// query := `
    //   SELECT *
    //   FROM public.subscriptions
    //   WHERE CURRENT_DATE <= end_date;
    //`
}

func (r *SubscriptionRepository) GetCurrentSubscription(id int) {}
