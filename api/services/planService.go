package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
)

type PlanService struct {
	PlanRepo *repository.PlanRepository
}

func NewPlanService(planRepo *repository.PlanRepository) *PlanService {
	return &PlanService{PlanRepo: planRepo}
}

func (s *PlanService) CreatePlan(d *model.PlanModel) (*model.PlanModel, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	_, err := s.PlanRepo.CreatePlan(d)

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

func (s *PlanService) GetASinglePlan(id int) (*model.PlanResponse, error) {
	plan, err := s.PlanRepo.GetSinglePlan(id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("no user found: %w", err)
	}
	return plan, nil
}

func (s *PlanService) UpdatePlan(d *model.PlanModel) error {
	if err := s.PlanRepo.EditPlan(d); err != nil {
		return err
	}
	return nil
}

func (s *PlanService) DeletePlan(id int) error {
	if err := s.PlanRepo.DeletePlan(id); err != nil {
		return err
	}
	return nil
}
