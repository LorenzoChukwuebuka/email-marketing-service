package repository

import (
	"email-marketing-service/api/v1/model"
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

	response := model.PlanResponse{
		UUID:                plan.UUID,
		PlanName:            plan.PlanName,
		Duration:            plan.Duration,
		Price:               plan.Price,
		NumberOfMailsPerDay: plan.NumberOfMailsPerDay,
		Details:             plan.Details,
		Status:              plan.Status,
		CreatedAt:           plan.CreatedAt,
	}

	if plan.UpdatedAt != nil {
		response.UpdatedAt = plan.UpdatedAt.Format(time.RFC3339)
	} else {
		response.UpdatedAt = ""
	}

	if plan.DeletedAt != nil {
		response.DeletedAt = plan.DeletedAt.Format(time.RFC3339)
	} else {
		response.DeletedAt = ""
	}

	return response
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
	existingPlan := model.Plan{}
	if err := r.DB.Model(&model.Plan{}).Where("uuid = ?", data.UUID).First(&existingPlan).Error; err != nil {
		return fmt.Errorf("failed to find plan for editing: %w", err)
	}

	// Update the existing plan with the new data
	existingPlan.PlanName = data.PlanName
	existingPlan.Duration = data.Duration
	existingPlan.Price = data.Price
	existingPlan.NumberOfMailsPerDay = data.NumberOfMailsPerDay
	existingPlan.Details = data.Details
	existingPlan.Status = data.Status

	// Save the changes to the database
	if err := r.DB.Save(&existingPlan).Error; err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}

func (r *PlanRepository) DeletePlan(id string) error {

	var existingPlan model.Plan
	if err := r.DB.Where("uuid = ?", id).First(&existingPlan).Error; err != nil {
		return fmt.Errorf("failed to find plan for deletion: %w", err)
	}

	// Soft delete by marking the plan as deleted
	if err := r.DB.Delete(&existingPlan).Error; err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	return nil
}
