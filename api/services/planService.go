package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
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

	return nil, nil
}

func (s *PlanService) GetAllPlans() ([]*model.PlanResponse, error) {
	return nil, nil
}

func (s *PlanService) GetASinglePlan(id int) (*model.PlanResponse, error) {
	return nil, nil
}

func (s *PlanService) UpdatePlan(id int) error {
	return nil
}

func (s *PlanService) DeletePlan(id int) error {
	return nil
}
