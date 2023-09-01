package controllers

import (
	"email-marketing-service/api/services"
	"net/http"
)

type PlanController struct {
	PlanService *services.PlanService
}

func NewPlanController(planService *services.PlanService) *PlanController {
	return &PlanController{
		PlanService: planService,
	}
}

func (c *PlanController) CreatePlan(w http.ResponseWriter, r *http.Request)    {}
func (c *PlanController) GetAllPlans(w http.ResponseWriter, r *http.Request)   {}
func (c *PlanController) GetSinglePlan(w http.ResponseWriter, r *http.Request) {}
func (c *PlanController) UpdatePlan(w http.ResponseWriter, r *http.Request)    {}
func (c *PlanController) DeletePlan(w http.ResponseWriter, r *http.Request)    {}
