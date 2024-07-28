package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPlanExists   = fmt.Errorf("plan already exists")
	ErrNoPlansFound = fmt.Errorf("no plans found")
	ErrNoPlanFound  = fmt.Errorf("no plan found")
)

type PlanService struct {
	planRepo *repository.PlanRepository
}

func NewPlanService(planRepo *repository.PlanRepository) *PlanService {
	return &PlanService{planRepo: planRepo}
}

func (s *PlanService) CreatePlan(d *dto.Plan) (*model.Plan, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid plan data: %w", err)
	}

	planModel := &model.Plan{
		UUID:                uuid.New().String(),
		PlanName:            d.PlanName,
		Duration:            d.Duration,
		Details:             d.Details,
		NumberOfMailsPerDay: d.NumberOfMailsPerDay,
		Price:               d.Price,
		Status:              model.PlanStatus(d.Status),
	}

	// Create slice of PlanFeature, but don't set PlanID yet
	features := make([]model.PlanFeature, len(d.Features))
	for i, feature := range d.Features {
		features[i] = model.PlanFeature{
			UUID:        uuid.New().String(),
			Name:        feature.Name,
			Identifier:  feature.Identifier,
			CountLimit:  feature.CountLimit,
			SizeLimit:   feature.SizeLimit,
			IsActive:    feature.IsActive,
			Description: feature.Description,
		}
	}

	exists, err := s.planRepo.PlanExistsByName(d.PlanName)
	if err != nil {
		return nil, fmt.Errorf("error checking plan existence: %w", err)
	}
	if exists {
		return nil, ErrPlanExists
	}

	// Pass both the plan and features to the repository
	createdPlan, err := s.planRepo.CreatePlan(planModel, features)
	if err != nil {
		return nil, fmt.Errorf("error creating plan: %w", err)
	}

	return createdPlan, nil
}

func (s *PlanService) GetAllPlans() ([]model.PlanResponse, error) {
	plans, err := s.planRepo.GetAllPlans()
	if err != nil {
		return nil, fmt.Errorf("error fetching plans: %w", err)
	}
	if len(plans) == 0 {
		return nil, ErrNoPlansFound
	}
	return plans, nil
}

func (s *PlanService) GetASinglePlan(id string) (model.PlanResponse, error) {
	plan, err := s.planRepo.GetSinglePlan(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PlanResponse{}, ErrNoPlanFound
		}
		return model.PlanResponse{}, fmt.Errorf("error fetching plan: %w", err)
	}
	return plan, nil
}

func (s *PlanService) UpdatePlan(d *dto.EditPlan) error {
	planModel := &model.Plan{
		UUID:                d.UUID,
		PlanName:            d.PlanName,
		Duration:            d.Duration,
		Details:             d.Details,
		NumberOfMailsPerDay: d.NumberOfMailsPerDay,
		Price:               d.Price,
		Status:              model.PlanStatus(d.Status),
		Features:            make([]model.PlanFeature, len(d.Features)),
	}

	for i, feature := range d.Features {
		planModel.Features[i] = model.PlanFeature{
			UUID:        uuid.New().String(),
			Name:        feature.Name,
			Identifier:  feature.Identifier,
			CountLimit:  feature.CountLimit,
			SizeLimit:   feature.SizeLimit,
			IsActive:    feature.IsActive,
			Description: feature.Description,
		}
	}

	if err := s.planRepo.EditPlan(planModel); err != nil {
		return fmt.Errorf("error updating plan: %w", err)
	}
	return nil
}

func (s *PlanService) DeletePlan(id string) error {
	if err := s.planRepo.DeletePlan(id); err != nil {
		return fmt.Errorf("error deleting plan: %w", err)
	}
	return nil
}
