package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"time"
)

type BillingRepository struct {
	DB *sql.DB
}

func NewBillingRepository(db *sql.DB) *BillingRepository {
	return &BillingRepository{DB: db}
}

var SetTime = func(field sql.NullTime, target *string) {
	if field.Valid {
		*target = field.Time.Format(time.RFC3339Nano)
	}
}

func (r *BillingRepository) CreateBilling(d *model.BillingModel) (*model.BillingModel, error) {
	query := "INSERT INTO billing (uuid,  user_id, amount_paid, plan_id, duration,  expiry_date, reference,paymentMethod,transaction_id, status, created_at) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8,$9,$10,$11 ) RETURNING id"

	stmt, err := r.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var insertedID int

	err = stmt.QueryRow(
		d.UUID,
		d.UserId,
		d.AmountPaid,
		d.PlanId,
		d.Duration,
		d.ExpiryDate,
		d.Reference,
		d.PaymentMethod,
		d.TransactionId,
		d.Status,
		d.CreatedAt,
	).Scan(&insertedID)

	d.Id = insertedID
	if err != nil {
		return nil, err
	}

	return d, nil
}
