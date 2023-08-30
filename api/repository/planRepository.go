package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
)

type PlanRepository struct {
	DB *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{
		DB: db,
	}
}

func (r *PlanRepository) CreatePlan() error {
	return nil
}

func (r *PlanRepository) GetAllPlans() ([]*model.PlanModel, error) {
	return nil, nil
}

func (r *PlanRepository) GetSinglePlan(id int) (*model.PlanModel, error) {
	return nil, nil
}

func (r *PlanRepository) EditPlan(data *model.PlanModel) error {
	return nil
}

func (r *PlanRepository) DeletePlan(id int) error {
	return nil
}
