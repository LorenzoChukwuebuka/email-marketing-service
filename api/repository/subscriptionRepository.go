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
	query := "INSERT INTO subscriptions (uuid,user_id, plan_id, payment_id, start_date, end_date, expired, created_at,updated_at,transaction_id,cancelled,date_cancelled) VALUES ($1, $2, $3, $4, $5, $6, $7,$8,$9,$10,$11,$12) RETURNING id"

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
		d.UpdatedAt,
		d.TransactionId,
		d.Cancelled,
		d.DateCancelled,
	).Scan(&insertedID)

	if err != nil {
		return err
	}

	d.Id = insertedID
	return nil // Return nil here on success
}

func (r *SubscriptionRepository) GetAllSubscriptions() ([]model.SubscriptionResponseModel, error) {
	query := `
        SELECT 
            s.id AS subscription_id,
            s.uuid AS subscription_uuid,
            s.user_id,
			u.id AS userId,
            u.uuid AS user_uuid,
            u.firstname,
            u.middlename,
            u.lastname,
            u.username,
            u.email,
            u.password,
            u.verified AS user_verified,
            u.created_at AS user_created_at,
            u.verified_at AS user_verified_at,
            u.updated_at AS user_updated_at,
            u.deleted_at AS user_deleted_at,
            s.plan_id,
			p.id AS planId,
            p.uuid AS plan_uuid,
            p.planname,
            p.duration,
            p.price,
            p.number_of_emails_per_day,
            p.details AS plan_details,
            p.status AS plan_status,
            p.created_at AS plan_created_at,
            p.updated_at AS plan_updated_at,
            p.deleted_at AS plan_deleted_at,
            s.payment_id,
            s.start_date,
            s.end_date,
            s.expired,
            s.transaction_id,
            s.created_at AS subscription_created_at,
            s.updated_at AS subscription_updated_at,
            s.cancelled AS subscription_cancelled,
            s.date_cancelled AS subscription_date_cancelled,
			b.amount_paid,
            b.expiry_date AS billing_expiry_date,
            b.reference,
            b.transaction_id AS billing_transaction_id,
            b.paymentmethod,
            b.status AS billing_status,
            b.created_at AS billing_created_at,
            b.updated_at AS billing_updated_at,
            b.deleted_at AS billing_deleted_at
        FROM 
            subscriptions s
        JOIN
            users u ON s.user_id = u.id
        JOIN
            plans p ON s.plan_id = p.id
	    JOIN
            billing b ON s.payment_id = b.id;
    `
	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var subscriptions []model.SubscriptionResponseModel

	for rows.Next() {
		var subscription model.SubscriptionResponseModel
		var updatedAt, dateCancelled, userVerifiedAt, userUpdatedAt, userDeletedAt, planUpdatedAt, planDeletedAt, billingUpdatedAt, billingDeletedAt sql.NullTime

		err := rows.Scan(
			&subscription.Id,
			&subscription.UUID,
			&subscription.UserId,
			&subscription.User.ID,
			&subscription.User.UUID,
			&subscription.User.FirstName,
			&subscription.User.MiddleName,
			&subscription.User.LastName,
			&subscription.User.UserName,
			&subscription.User.Email,
			&subscription.User.Password,
			&subscription.User.Verified,
			&subscription.User.CreatedAt,
			&userVerifiedAt,
			&userUpdatedAt,
			&userDeletedAt,
			&subscription.PlanId,
			&subscription.Plan.Id,
			&subscription.Plan.UUID,
			&subscription.Plan.PlanName,
			&subscription.Plan.Duration,
			&subscription.Plan.Price,
			&subscription.Plan.NumberOfMailsPerDay,
			&subscription.Plan.Details,
			&subscription.Plan.Status,
			&subscription.Plan.CreatedAt,
			&planUpdatedAt,
			&planDeletedAt,
			&subscription.PaymentId,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.Expired,
			&subscription.TransactionId,
			&subscription.CreatedAt,
			&updatedAt,
			&subscription.Cancelled,
			&dateCancelled,
			&subscription.Billing.AmountPaid,
			//&subscription.Billing.Duration,
			&subscription.Billing.ExpiryDate,
			&subscription.Billing.Reference,
			&subscription.Billing.TransactionId,
			&subscription.Billing.PaymentMethod,
			&subscription.Billing.Status,
			&subscription.Billing.CreatedAt,
			&billingUpdatedAt,
			&billingDeletedAt,
		)

		if err != nil {
			return nil, err
		}

		SetTime(updatedAt, &subscription.UpdatedAt)
		SetTime(dateCancelled, &subscription.DateCancelled)
		SetTime(userVerifiedAt, &subscription.User.VerifiedAt)
		SetTime(userUpdatedAt, &subscription.User.UpdatedAt)
		SetTime(userDeletedAt, &subscription.User.DeletedAt)
		SetTime(planDeletedAt, &subscription.Plan.DeletedAt)
		SetTime(planUpdatedAt, &subscription.Plan.UpdatedAt)
		SetTime(billingUpdatedAt, &subscription.Billing.UpdatedAt)
		SetTime(billingDeletedAt, &subscription.Billing.DeletedAt)

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
			&subscription.Cancelled,
			&subscription.DateCancelled,
		)

		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, err
}

func (r *SubscriptionRepository) UpdateExpiredSubscription(id int) error {
	query := "UPDATE subscriptions SET expired = true WHERE id = $1"

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
