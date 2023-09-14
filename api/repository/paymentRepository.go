package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"time"
)

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) CreatePayment(d *model.PaymentModel) (*model.PaymentModel, error) {
	query := "INSERT INTO payments (  user_id, amount_paid, plan_id, duration,  expiry_date, reference, status, created_at) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 ) RETURNING id"

	stmt, err := r.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var insertedID int

	err = stmt.QueryRow(
		d.UserId,
		d.AmountPaid,
		d.PlanId,
		d.Duration,
		d.ExpiryDate,
		d.Reference,
		d.Status,
		d.CreatedAt,
	).Scan(&insertedID)

	d.Id = insertedID
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (r *PaymentRepository) GetSinglePayment(id int,userId int) (*model.PaymentResponse, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			p.amount_paid,
			p.plan_id,
			p.duration,
			p.expiry_date,
			p.reference,
			p.status,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			u.id AS "user.id",
			u.uuid AS "user.uuid",
			u.first_name AS "user.firstname",
			u.middle_name AS "user.middlename",
			u.last_name AS "user.lastname",
			u.username AS "user.username",
			u.email AS "user.email",
			u.password AS "user.password",
			u.verified AS "user.verified",
			u.created_at AS "user.created_at",
			u.verified_at AS "user.verified_at",
			u.updated_at AS "user.updated_at",
			u.deleted_at AS "user.deleted_at",
			pl.id AS "plan.id",
			pl.planname AS "plan.planname",
			pl.duration AS "plan.duration",
			pl.price AS "plan.price",
			pl.details AS "plan.details",
			pl.status AS "plan.status",
			pl.created_at AS "plan.created_at",
			pl.updated_at AS "plan.updated_at",
			pl.deleted_at AS "plan.deleted_at"
		FROM
			payments p
		JOIN
			users u ON p.user_id = u.id
		JOIN
			plans pl ON p.plan_id = pl.id
		WHERE
			p.id = $1
		AND WHERE 
		   u.id = $2	;
	`

	// Prepare the SQL statement
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement
	row := stmt.QueryRow(id)

	// Scan the result into a PaymentResponse struct
	var payment model.PaymentResponse

	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	var userverified sql.NullTime
	var userdeleted sql.NullTime
	var userupdated sql.NullTime

	var planupated sql.NullTime
	var plandeleted sql.NullTime

	err = row.Scan(
		&payment.Id,
		&payment.UserId,
		&payment.AmountPaid,
		&payment.PlanId,
		&payment.Duration,
		&payment.ExpiryDate,
		&payment.Reference,
		&payment.Status,
		&payment.CreatedAt,
		&updatedAt,
		&deletedAt,
		&payment.User.ID,
		&payment.User.UUID,
		&payment.User.FirstName,
		&payment.User.MiddleName,
		&payment.User.LastName,
		&payment.User.UserName,
		&payment.User.Email,
		&payment.User.Password,
		&payment.User.Verified,
		&payment.User.CreatedAt,
		&userverified,
		&userupdated,
		&userdeleted,
		&payment.Plan.Id,
		&payment.Plan.PlanName,
		&payment.Plan.Duration,
		&payment.Plan.Price,
		&payment.Plan.Details,
		&payment.Plan.Status,
		&payment.Plan.CreatedAt,
		&planupated,
		&plandeleted,
	)
	if err != nil {
		return nil, err
	}

	// Set UpdatedAt based on the presence of updatedAt or deletedAt
	if deletedAt.Valid {
		payment.UpdatedAt = deletedAt.Time.Format(time.RFC3339Nano)
	}
	if updatedAt.Valid {
		payment.UpdatedAt = updatedAt.Time.Format(time.RFC3339Nano)
	}

	if userverified.Valid {
		payment.User.VerifiedAt = userverified.Time.Format(time.RFC3339Nano)
	}

	if userupdated.Valid {
		payment.User.UpdatedAt = userupdated.Time.Format(time.RFC3339Nano)
	}

	if userdeleted.Valid {
		payment.User.DeletedAt = userdeleted.Time.Format(time.RFC3339Nano)
	}

	if plandeleted.Valid {
		payment.Plan.DeletedAt = plandeleted.Time.Format(time.RFC3339Nano)
	}

	if planupated.Valid {
		payment.Plan.UpdatedAt = plandeleted.Time.Format(time.RFC3339Nano)
	}

	return &payment, nil
}

func (r *PaymentRepository) GetAllPayments(userId int) ([]model.PaymentResponse, error) {
	query := `SELECT
			p.id,
			p.user_id,
			p.amount_paid,
			p.plan_id,
			p.duration,
			p.expiry_date,
			p.reference,
			p.status,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			u.id AS "user.id",
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
			pl.id AS "plan.id",
			pl.planname AS "plan.planname",
			pl.duration AS "plan.duration",
			pl.price AS "plan.price",
			pl.details AS "plan.details",
			pl.status AS "plan.status",
			pl.created_at AS "plan.created_at",
			pl.updated_at AS "plan.updated_at",
			pl.deleted_at AS "plan.deleted_at"
		FROM
			payments p
		JOIN
			users u ON p.user_id = u.id
		JOIN
			plans pl ON p.plan_id = pl.id

			where u.id = $1
		`

	// Prepare the SQL statement
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement
	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and map to PaymentResponse struct
	var payments []model.PaymentResponse
	for rows.Next() {
		var payment model.PaymentResponse
		var updatedAt sql.NullTime
		var deletedAt sql.NullTime

		var userverified sql.NullTime
		var userdeleted sql.NullTime
		var userupdated sql.NullTime

		var planupated sql.NullTime
		var plandeleted sql.NullTime

		err := rows.Scan(
			&payment.Id,
			&payment.UserId,
			&payment.AmountPaid,
			&payment.PlanId,
			&payment.Duration,
			&payment.ExpiryDate,
			&payment.Reference,
			&payment.Status,
			&payment.CreatedAt,
			&updatedAt,
			&deletedAt,
			&payment.User.ID,
			&payment.User.UUID,
			&payment.User.FirstName,
			&payment.User.MiddleName,
			&payment.User.LastName,
			&payment.User.UserName,
			&payment.User.Email,
			&payment.User.Password,
			&payment.User.Verified,
			&payment.User.CreatedAt,
			&userverified,
			&userupdated,
			&userdeleted,
			&payment.Plan.Id,
			&payment.Plan.PlanName,
			&payment.Plan.Duration,
			&payment.Plan.Price,
			&payment.Plan.Details,
			&payment.Plan.Status,
			&payment.Plan.CreatedAt,
			&planupated,
			&plandeleted,
		)
		if err != nil {
			return nil, err
		}

		// Set UpdatedAt based on the presence of updatedAt or deletedAt
		if deletedAt.Valid {
			payment.UpdatedAt = deletedAt.Time.Format(time.RFC3339Nano)
		}
		if updatedAt.Valid {
			payment.UpdatedAt = updatedAt.Time.Format(time.RFC3339Nano)
		}

		if userverified.Valid {
			payment.User.VerifiedAt = userverified.Time.Format(time.RFC3339Nano)
		}

		if userupdated.Valid {
			payment.User.UpdatedAt = userupdated.Time.Format(time.RFC3339Nano)
		}

		if userdeleted.Valid {
			payment.User.DeletedAt = userdeleted.Time.Format(time.RFC3339Nano)
		}

		if plandeleted.Valid {
			payment.Plan.DeletedAt = plandeleted.Time.Format(time.RFC3339Nano)
		}

		if planupated.Valid {
			payment.Plan.UpdatedAt = plandeleted.Time.Format(time.RFC3339Nano)
		}

		payments = append(payments, payment)
	}

	return payments, nil
}
