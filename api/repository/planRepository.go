package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"time"
)

type PlanRepository struct {
	DB *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{
		DB: db,
	}
}

func (r *PlanRepository) CreatePlan(d *model.PlanModel) (*model.PlanModel, error) {
	query := "INSERT INTO plans(planname, duration, price, details, status, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var insertedID int
	err = stmt.QueryRow(d.PlanName, d.Duration, d.Price, d.Details, "active", time.Now()).Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	d.Id = insertedID
	return d, nil
}

func (r *PlanRepository) PlanExistsByName(planname string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM plans WHERE planname = $1)"
	var count bool
	err := r.DB.QueryRow(query, planname).Scan(&count)
	if err != nil {
		return false, err
	}
	return count, nil
}

func (r *PlanRepository) GetAllPlans() ([]model.PlanResponse, error) {
	query := "Select * from plans"

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plans []model.PlanResponse

	for rows.Next() {
		var plan model.PlanResponse
		var updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&plan.Id,
			&plan.PlanName,
			&plan.Duration,
			&plan.Price,
			&plan.Details,
			&plan.Status,
			&plan.CreatedAt,
			&updatedAt,
			&deletedAt,
		)

		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			plan.UpdatedAt = deletedAt.Time.Format(time.RFC3339Nano)
		}
		if updatedAt.Valid {
			plan.UpdatedAt = updatedAt.Time.Format(time.RFC3339Nano)
		}

		plans = append(plans, plan)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func (r *PlanRepository) GetSinglePlan(id int) (*model.PlanResponse, error) {
	query := "Select * FROM plans WHERE id = $1"

	row := r.DB.QueryRow(query, id)

	var plan model.PlanResponse
	var updatedAt, deletedAt sql.NullTime
	err := row.Scan(
		&plan.Id,
		&plan.PlanName,
		&plan.Duration,
		&plan.Price,
		&plan.Details,
		&plan.Status,
		&plan.CreatedAt,
		&updatedAt,
		&deletedAt,
	)

	if deletedAt.Valid {
		plan.UpdatedAt = deletedAt.Time.Format(time.RFC3339Nano)
	}
	if updatedAt.Valid {
		plan.UpdatedAt = updatedAt.Time.Format(time.RFC3339Nano)
	}

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) EditPlan(data *model.PlanModel) error {
	return nil
}

func (r *PlanRepository) DeletePlan(id int) error {
	return nil
}
