package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type PlanRepository struct {
	DB *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{
		DB: db,
	}
}

func (r *PlanRepository) createPlanResponse(plan model.Plan) model.PlanResponse {

	return model.PlanResponse{
		UUID:                plan.UUID,
		PlanName:            plan.PlanName,
		Duration:            plan.Duration,
		Price:               plan.Price,
		NumberOfMailsPerDay: plan.NumberOfMailsPerDay,
		Details:             plan.Details,
		Status:              plan.Status,
		CreatedAt:           plan.CreatedAt,
		UpdatedAt:           plan.UpdatedAt.Format(time.RFC3339),
		DeletedAt:           plan.DeletedAt.Format(time.RFC3339),
	}
}

func (r *PlanRepository) CreatePlan(d *model.Plan) (*model.Plan, error) {
	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert plan: %w", err)
	}

	return d, nil
}

func (r *PlanRepository) PlanExistsByName(planname string) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.Plan{}).Where("plan_name = ?", planname).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PlanRepository) GetAllPlans() ([]model.PlanResponse, error) {
	var plans []model.PlanResponse
	if err := r.DB.Model(&model.Plan{}).Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepository) GetSinglePlan(uuid string) (model.PlanResponse, error) {
	var plan model.Plan
	if err := r.DB.Model(&model.Plan{}).Where("uuid = ?", uuid).First(&plan).Error; err != nil {
		return model.PlanResponse{}, err
	}
	return r.createPlanResponse(plan), nil
}

func (r *PlanRepository) EditPlan(data *model.Plan) error {

	return nil
}

func (r *PlanRepository) DeletePlan(id string) error {

	return nil
}
