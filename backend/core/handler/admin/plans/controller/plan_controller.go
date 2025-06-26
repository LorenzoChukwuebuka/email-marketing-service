package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/plans/dto"
	"email-marketing-service/core/handler/admin/plans/service"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type PlanController struct {
	service *service.PlanService
}

func NewPlanController(service *service.PlanService) *PlanController {
	return &PlanController{
		service: service,
	}
}

func (c *PlanController) CreatePlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req *dto.Plan

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	createdPlan, err := c.service.CreatePlan(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusCreated, createdPlan)
}

func (c *PlanController) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	plans, err := c.service.GetAllPlansWithDetails(ctx)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, plans)
}

func (c *PlanController) GetPlanByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)

	planId := vars["planId"]

	planID, err := common.ParseUUIDMap(map[string]string{"planId": planId})
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	plan, err := c.service.GetPlanByID(ctx, planID["planId"])
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, plan)
}

func (c *PlanController) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)

	planId := vars["planId"]

	var req *dto.EditPlan

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	req.UUID = planId

	updatedPlan, err := c.service.UpdatePlan(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, updatedPlan)
}

func (c *PlanController) DeletePlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)

	planId := vars["planId"]

	planID, err := common.ParseUUIDMap(map[string]string{"planId": planId})
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	err = c.service.DeletePlan(ctx, planID["planId"])
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, nil)
}
