package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"net/http"
	"github.com/gorilla/mux"
)

type PlanController struct {
	PlanService *services.PlanService
}

func NewPlanController(planService *services.PlanService) *PlanController {
	return &PlanController{
		PlanService: planService,
	}
}

func (c *PlanController) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.Plan

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.PlanService.CreatePlan(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *PlanController) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	result, err := c.PlanService.GetAllPlans()

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *PlanController) GetSinglePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := c.PlanService.GetASinglePlan(id)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *PlanController) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var reqdata *dto.EditPlan

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.UUID = id

	err := c.PlanService.UpdatePlan(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "plan edited successfully")

}

func (c *PlanController) DeletePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	if err := c.PlanService.DeletePlan(id); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "plan deleted successfully")
}
