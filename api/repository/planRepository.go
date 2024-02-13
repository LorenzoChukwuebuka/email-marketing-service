package repository

import (
	"email-marketing-service/api/model"
	"gorm.io/gorm"
)

type PlanRepository struct {
	DB *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{
		DB: db,
	}
}

func (r *PlanRepository) CreatePlan(d *model.PlanModel) (*model.PlanModel, error) {
	return nil, nil
}

func (r *PlanRepository) PlanExistsByName(planname string) (bool, error) {
	return false, nil
}

func (r *PlanRepository) GetAllPlans() ([]model.PlanResponse, error) {
	return nil, nil
}

func (r *PlanRepository) GetSinglePlan(id string) (*model.PlanResponse, error) {
	return nil, nil
}

func (r *PlanRepository) EditPlan(data *model.PlanModel) error {

	return nil
}

func (r *PlanRepository) DeletePlan(id string) error {

	return nil
}
