package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	"email-marketing-service/core/handler/admin/email-templates/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type AdminUserTemplateController struct {
	service *services.AdminUserTemplatesService
}

func NewAdminUserTemplateController(service *services.AdminUserTemplatesService) *AdminUserTemplateController {
	return &AdminUserTemplateController{
		service: service,
	}
}

func (c *AdminUserTemplateController) GetUserTemplates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	userId := vars["userId"]
	tempType := vars["type"]

	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	req := &dto.AdminFetchTemplateDTO{
		UserID:      userId,
		SearchQuery: search,
		Offset:      offset,
		Limit:       limit,
		Type:        tempType,
	}

	templates, err := c.service.GetUserTemplates(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to fetch user templates: %v", err), nil)
		return
	}

	helper.SuccessResponse(w, 200, templates)
}

func (c *AdminUserTemplateController) GetSingleTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	templateId := vars["templateId"]
	userId := vars["userId"]

	req := &dto.AdminFetchTemplateDTO{
		TemplateId: templateId,
		UserID:     userId,
	}

	template, err := c.service.GetUserTemplateById(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to fetch template: %v", err), nil)
		return
	}

	helper.SuccessResponse(w, 200, template)
}


