package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	"email-marketing-service/core/handler/admin/email-templates/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type AdminTemplateController struct {
	service *services.AdminTemplatesService
}

func NewAdminTemplateController(service *services.AdminTemplatesService) *AdminTemplateController {
	return &AdminTemplateController{
		service: service,
	}
}

func (c *AdminTemplateController) CreateGalleryTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.AdminTemplateDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	result, err := c.service.CreateTemplate(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to create template: %v", err), nil)
		return
	}

	helper.SuccessResponse(w, 201, result)
}

func (c *AdminTemplateController) GetTemplateById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	templateId := vars["templateId"]

	result, err := c.service.GetTemplate(ctx, templateId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminTemplateController) GetTemplatesByType(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	tempType := vars["type"]

	page, pageSize, search, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	if err != nil {
		helper.ErrorResponse(w, err, common.ErrGenericError)
		return
	}

	req := &dto.AdminFetchGalleryTemplatesDTO{
		Limit:  limit,
		Offset: offset,
		Type:   tempType,
		Search: search,
	}

	result, err := c.service.GetTemplatesByType(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to fetch template: %v", err), nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminTemplateController) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.AdminTemplateDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	vars := mux.Vars(r)
	tempId := vars["templateId"]

	req.TemplateID = tempId

	result, err := c.service.UpdateTemplate(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to update template:%w", err), nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminTemplateController) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	tempId := vars["templateId"]

	if err := c.service.DeleteTemplate(ctx, tempId); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to delete template:%w", err), nil)
		return
	}

	helper.SuccessResponse(w, 200, "template deleted successfully")
}
