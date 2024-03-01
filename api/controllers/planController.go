package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	var reqdata *model.Plan

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

	plan_id, _ := strconv.Atoi(id)

	var reqdata *model.Plan

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.ID = plan_id

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
