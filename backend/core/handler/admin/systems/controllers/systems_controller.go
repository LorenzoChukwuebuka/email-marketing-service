package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/systems/dto"
	"email-marketing-service/core/handler/admin/systems/services"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type SystemsController struct {
	SystemsService *services.Service
}

func NewAdminSystemsController(systemsSvc *services.Service) *SystemsController {
	return &SystemsController{
		SystemsService: systemsSvc,
	}
}

func (c *SystemsController) CreateRecords(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	var reqdata *dto.SystemsDTO

	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	result, err := c.SystemsService.GenerateAndSaveSMTPCredentials(ctx, reqdata.Domain)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 201, result)
}

func (c *SystemsController) GetDNSRecords(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)

	domainName := vars["domain"]

	result, err := c.SystemsService.GetDNSRecords(ctx, domainName)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *SystemsController) DeleteDNSRecords(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	domainName := vars["domain"]
	err := c.SystemsService.DeleteDNSRecords(ctx, domainName)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "deleted successfully")
}

func (c *SystemsController) ReadAppLogs(w http.ResponseWriter, r *http.Request) {
	result, err := c.SystemsService.ReadAppLogs()
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *SystemsController) ReadRequestLogs(w http.ResponseWriter, r *http.Request) {
	result, err := c.SystemsService.ReadRequestLogs()
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}
