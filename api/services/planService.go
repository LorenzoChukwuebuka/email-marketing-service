package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/google/uuid"
)

var EmptyPlanResponse = model.PlanResponse{} // Zero-initialized instance

type PlanService struct {
	PlanRepo *repository.PlanRepository
}

func NewPlanService(planRepo *repository.PlanRepository) *PlanService {
	return &PlanService{PlanRepo: planRepo}
}

func (s *PlanService) CreatePlan(d *model.Plan) (*model.Plan, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	d.UUID = uuid.New().String()

	planExists, err := s.PlanRepo.PlanExistsByName(d.PlanName)

	if err != nil {
		return nil, err
	}

	if planExists {
		return nil, fmt.Errorf("plan already exists")
	}

	_, err = s.PlanRepo.CreatePlan(d)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *PlanService) GetAllPlans() ([]model.PlanResponse, error) {
	plans, err := s.PlanRepo.GetAllPlans()

	if err != nil {
		return nil, err
	}

	if len(plans) == 0 {
		return nil, fmt.Errorf("no user found: %w", err)
	}
	return plans, nil
}

func (s *PlanService) GetASinglePlan(id string) (model.PlanResponse, error) {
	plan, err := s.PlanRepo.GetSinglePlan(id)
	if err != nil {
		return EmptyPlanResponse, err
	}
	if plan == EmptyPlanResponse {
		return EmptyPlanResponse, fmt.Errorf("no record found")
	}
	return plan, nil
}

func (s *PlanService) UpdatePlan(d *model.Plan) error {
	if err := s.PlanRepo.EditPlan(d); err != nil {
		return err
	}
	return nil
}

func (s *PlanService) DeletePlan(id string) error {
	if err := s.PlanRepo.DeletePlan(id); err != nil {
		return err
	}
	return nil
}
