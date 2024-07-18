package services

import (
	"fmt"

	"github.com/google/uuid"

	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
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
		Status:              d.Status,
	}

	exists, err := s.planRepo.PlanExistsByName(d.PlanName)
	if err != nil {
		return nil, fmt.Errorf("error checking plan existence: %w", err)
	}
	if exists {
		return nil, ErrPlanExists
	}

	_, err = s.planRepo.CreatePlan(planModel)
	if err != nil {
		return nil, fmt.Errorf("error creating plan: %w", err)
	}

	return planModel, nil
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
		return model.PlanResponse{}, fmt.Errorf("error fetching plan: %w", err)
	}
	if plan == (model.PlanResponse{}) {
		return model.PlanResponse{}, ErrNoPlanFound
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
		Status:              d.Status,
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
