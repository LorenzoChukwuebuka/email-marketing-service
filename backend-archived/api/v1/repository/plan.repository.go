package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"strconv"
	"time"

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

func (r *PlanRepository) createPlanResponse(plan model.Plan) model.PlanResponse {
	response := model.PlanResponse{
		ID:                  plan.ID,
		UUID:                plan.UUID,
		PlanName:            plan.PlanName,
		Duration:            plan.Duration,
		Price:               plan.Price,
		NumberOfMailsPerDay: plan.NumberOfMailsPerDay,
		Details:             plan.Details,
		Status:              plan.Status,
		Features:            plan.Features,
		MailingLimit:        plan.MailingLimit,
		CreatedAt:           plan.CreatedAt.String(),
		UpdatedAt:           plan.UpdatedAt.String(),
	}

	if plan.DeletedAt.Valid {
		htime := plan.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &htime
	}

	return response
}

func (r *PlanRepository) CreatePlan(plan *model.Plan, features []model.PlanFeature) (*model.Plan, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Create the plan
		if err := tx.Create(plan).Error; err != nil {
			return fmt.Errorf("failed to insert plan: %w", err)
		}

		limitAmount, err := strconv.Atoi(plan.NumberOfMailsPerDay)
		if err != nil {
			return fmt.Errorf("failed to parse NumberOfMailsPerDay: %w", err)
		}

		mailingLimit := model.MailingLimit{
			PlanID:      plan.ID,
			LimitAmount: limitAmount,
			LimitPeriod: func() string {
				if plan.PlanName == "free" {
					return "day"
				}
				return "month"
			}(),
		}

		if err := tx.Create(&mailingLimit).Error; err != nil {
			return fmt.Errorf("failed to insert mailing limit: %w", err)
		}

		// Create features (existing code)
		for i := range features {
			features[i].PlanID = plan.ID
			if err := tx.Create(&features[i]).Error; err != nil {
				return fmt.Errorf("failed to insert plan feature: %w", err)
			}
		}

		// Set the features on the plan
		plan.Features = features
		plan.MailingLimit = mailingLimit

		return nil
	})

	if err != nil {
		return nil, err
	}

	return plan, nil

}

func (r *PlanRepository) PlanExistsByName(planname string) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.Plan{}).Where("plan_name = ?", planname).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PlanRepository) GetAllPlans() ([]model.PlanResponse, error) {
	var plans []model.Plan
	if err := r.DB.Model(&model.Plan{}).Preload("Features").Preload("MailingLimit").Find(&plans).Error; err != nil {
		return nil, err
	}

	var planResponses []model.PlanResponse
	for _, plan := range plans {
		planResponses = append(planResponses, r.createPlanResponse(plan))
	}

	return planResponses, nil
}

func (r *PlanRepository) GetSinglePlan(uuid string) (model.PlanResponse, error) {
	var plan model.Plan
	if err := r.DB.Model(&model.Plan{}).Preload("Features").Preload("MailingLimit").Where("uuid = ?", uuid).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.PlanResponse{}, gorm.ErrRecordNotFound
		}
		return model.PlanResponse{}, err
	}
	return r.createPlanResponse(plan), nil
}

func (r *PlanRepository) GetPlanById(id uint) (model.PlanResponse, error) {
	var plan model.Plan
	if err := r.DB.Model(&model.Plan{}).Preload("Features").Preload("MailingLimit").Where("id = ?", id).First(&plan).Error; err != nil {
		return model.PlanResponse{}, err
	}
	return r.createPlanResponse(plan), nil
}

func (r *PlanRepository) EditPlan(data *model.Plan) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		existingPlan := model.Plan{}
		if err := tx.Preload("Features").Preload("MailingLimit").Where("uuid = ?", data.UUID).First(&existingPlan).Error; err != nil {
			return fmt.Errorf("failed to find plan for editing: %w", err)
		}

		// Update the existing plan with the new data
		existingPlan.PlanName = data.PlanName
		existingPlan.Duration = data.Duration
		existingPlan.Price = data.Price
		existingPlan.NumberOfMailsPerDay = data.NumberOfMailsPerDay
		existingPlan.Details = data.Details
		existingPlan.Status = data.Status

		limitAmount, err := strconv.Atoi(data.NumberOfMailsPerDay)
		if err != nil {
			return fmt.Errorf("failed to parse NumberOfMailsPerDay: %w", err)
		}

		// Update MailingLimit
		existingPlan.MailingLimit.LimitAmount = limitAmount
		if err := tx.Save(&existingPlan.MailingLimit).Error; err != nil {
			return fmt.Errorf("failed to update mailing limit: %w", err)
		}

		// Update features (existing code)
		if err := tx.Where("plan_id = ?", existingPlan.ID).Delete(&model.PlanFeature{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing features: %w", err)
		}

		existingPlan.Features = data.Features
		for i := range existingPlan.Features {
			existingPlan.Features[i].PlanID = existingPlan.ID
		}

		// Save the changes to the database
		if err := tx.Save(&existingPlan).Error; err != nil {
			return fmt.Errorf("failed to update plan: %w", err)
		}

		return nil
	})
}

func (r *PlanRepository) DeletePlan(id string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var existingPlan model.Plan
		if err := tx.Where("uuid = ?", id).First(&existingPlan).Error; err != nil {
			return fmt.Errorf("failed to find plan for deletion: %w", err)
		}

		// Delete associated features
		if err := tx.Where("plan_id = ?", existingPlan.ID).Delete(&model.PlanFeature{}).Error; err != nil {
			return fmt.Errorf("failed to delete associated features: %w", err)
		}

		// Delete associated MailingLimit
		if err := tx.Where("plan_id = ?", existingPlan.ID).Delete(&model.MailingLimit{}).Error; err != nil {
			return fmt.Errorf("failed to delete associated mailing limit: %w", err)
		}

		// Soft delete the plan
		if err := tx.Delete(&existingPlan).Error; err != nil {
			return fmt.Errorf("failed to delete plan: %w", err)
		}

		return nil
	})
}
