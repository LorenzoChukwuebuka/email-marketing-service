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

func (r *BillingRepository) GetSingleBillingRecord(billingID string, userID int) (*model.BillingResponse, error) {
	query := `
			SELECT
			p.id,
			p.uuid,
			p.user_id AS "user.id",
			p.amount_paid,
			p.plan_id AS "plan.id",
			p.duration,
			p.expiry_date,
			p.reference,
			p.transaction_id,
			p.paymentMethod AS "payment_method",
			p.status,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			u.uuid AS "user.uuid",
			u.firstname AS "user.firstname",
			u.middlename AS "user.middlename",
			u.lastname AS "user.lastname",
			u.username AS "user.username",
			u.email AS "user.email",
			u.password AS "user.password",
			u.verified AS "user.verified",
			u.created_at AS "user.created_at",
			u.verified_at AS "user.verified_at",
			u.updated_at AS "user.updated_at",
			u.deleted_at AS "user.deleted_at",
			pl.uuid AS "plan.uuid",
			pl.planname AS "plan.planname",
			pl.duration AS "plan.duration",
			pl.price AS "plan.price",
			pl.number_of_emails_per_day,
			pl.details AS "plan.details",
			pl.status AS "plan.status",
			pl.created_at AS "plan.created_at",
			pl.updated_at AS "plan.updated_at",
			pl.deleted_at AS "plan.deleted_at"
		FROM
			billing p
		JOIN
			users u ON p.user_id = u.id
		JOIN
			plans pl ON p.plan_id = pl.id
		WHERE
			u.id = $1
		AND 
			p.uuid = $2
	`

	// Prepare the SQL statement
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement
	row := stmt.QueryRow(userID, billingID)

	// Scan the result into a PaymentResponse struct
	var payment model.BillingResponse

	var updatedAt, deletedAt, userVerifiedAt, userUpdatedAt, userDeletedAt, planUpdatedAt, planDeletedAt sql.NullTime

	err = row.Scan(
		&payment.Id,
		&payment.UUID,
		&payment.UserId,
		&payment.AmountPaid,
		&payment.PlanId,
		&payment.Duration,
		&payment.ExpiryDate,
		&payment.Reference,
		&payment.TransactionId,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.CreatedAt,
		&updatedAt,
		&deletedAt,
		&payment.User.UUID,
		&payment.User.FirstName,
		&payment.User.MiddleName,
		&payment.User.LastName,
		&payment.User.UserName,
		&payment.User.Email,
		&payment.User.Password,
		&payment.User.Verified,
		&payment.User.CreatedAt,
		&userVerifiedAt,
		&userUpdatedAt,
		&userDeletedAt,
		&payment.Plan.UUID,
		&payment.Plan.PlanName,
		&payment.Plan.Duration,
		&payment.Plan.Price,
		&payment.Plan.NumberOfMailsPerDay,
		&payment.Plan.Details,
		&payment.Plan.Status,
		&payment.Plan.CreatedAt,
		&planUpdatedAt,
		&planDeletedAt,
	)
	if err != nil {
		return nil, err
	}

	// Convert the sql.NullTime and sql.NullString to string
	SetTime(updatedAt, &payment.UpdatedAt)
	SetTime(deletedAt, &payment.DeletedAt)
	SetTime(userVerifiedAt, &payment.User.VerifiedAt)
	SetTime(userUpdatedAt, &payment.User.UpdatedAt)
	SetTime(userDeletedAt, &payment.User.DeletedAt)
	SetTime(planDeletedAt, &payment.Plan.DeletedAt)
	SetTime(planUpdatedAt, &payment.Plan.UpdatedAt)

	return &payment, nil
}

func (r *BillingRepository) GetAllPayments(userId int, page int) ([]model.BillingResponse, error) {

	// Assuming a fixed page size of 20
	pageSize := 20

	// Calculate the offset based on the page number and fixed page size
	offset := (page - 1) * pageSize

	query := `
	SELECT
			p.id,
			p.uuid,
			p.user_id AS "user.id",
			p.amount_paid,
			p.plan_id AS "plan.id",
			p.duration,
			p.expiry_date,
			p.reference,
			p.transaction_id,
			p.paymentMethod AS "payment_method",
			p.status,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			u.uuid AS "user.uuid",
			u.firstname AS "user.firstname",
			u.middlename AS "user.middlename",
			u.lastname AS "user.lastname",
			u.username AS "user.username",
			u.email AS "user.email",
			u.password AS "user.password",
			u.verified AS "user.verified",
			u.created_at AS "user.created_at",
			u.verified_at AS "user.verified_at",
			u.updated_at AS "user.updated_at",
			u.deleted_at AS "user.deleted_at",
			pl.uuid AS "plan.uuid",
			pl.planname AS "plan.planname",
			pl.duration AS "plan.duration",
			pl.price AS "plan.price",
			pl.number_of_emails_per_day,
			pl.details AS "plan.details",
			pl.status AS "plan.status",
			pl.created_at AS "plan.created_at",
			pl.updated_at AS "plan.updated_at",
			pl.deleted_at AS "plan.deleted_at"
		FROM
			billing p
		JOIN
			users u ON p.user_id = u.id
		JOIN
			plans pl ON p.plan_id = pl.id
		WHERE
			u.id = $1
		ORDER BY 
		    p.created_at DESC  -- You can adjust the ORDER BY clause as needed
		LIMIT $2 OFFSET $3
`

	// Prepare the SQL statement
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement
	rows, err := stmt.Query(userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and map to PaymentResponse struct
	var payments []model.BillingResponse
	for rows.Next() {
		var payment model.BillingResponse
		var updatedAt, deletedAt, userVerifiedAt, userUpdatedAt, userDeletedAt, planUpdatedAt, planDeletedAt sql.NullTime

		err := rows.Scan(
			&payment.Id,
			&payment.UUID,
			&payment.UserId,
			&payment.AmountPaid,
			&payment.PlanId,
			&payment.Duration,
			&payment.ExpiryDate,
			&payment.Reference,
			&payment.TransactionId,
			&payment.PaymentMethod,
			&payment.Status,
			&payment.CreatedAt,
			&updatedAt,
			&deletedAt,
			&payment.User.UUID,
			&payment.User.FirstName,
			&payment.User.MiddleName,
			&payment.User.LastName,
			&payment.User.UserName,
			&payment.User.Email,
			&payment.User.Password,
			&payment.User.Verified,
			&payment.User.CreatedAt,
			&userVerifiedAt,
			&userUpdatedAt,
			&userDeletedAt,
			&payment.Plan.UUID,
			&payment.Plan.PlanName,
			&payment.Plan.Duration,
			&payment.Plan.Price,
			&payment.Plan.NumberOfMailsPerDay,
			&payment.Plan.Details,
			&payment.Plan.Status,
			&payment.Plan.CreatedAt,
			&planUpdatedAt,
			&planDeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert the sql.NullTime and sql.NullString to string
		SetTime(updatedAt, &payment.UpdatedAt)
		SetTime(deletedAt, &payment.DeletedAt)
		SetTime(userVerifiedAt, &payment.User.VerifiedAt)
		SetTime(userUpdatedAt, &payment.User.UpdatedAt)
		SetTime(userDeletedAt, &payment.User.DeletedAt)
		SetTime(planDeletedAt, &payment.Plan.DeletedAt)
		SetTime(planUpdatedAt, &payment.Plan.UpdatedAt)

		payments = append(payments, payment)
	}

	return payments, nil
}
