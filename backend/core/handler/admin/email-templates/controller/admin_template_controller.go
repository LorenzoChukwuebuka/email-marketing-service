package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	"email-marketing-service/core/handler/admin/email-templates/services"
	"email-marketing-service/internal/helper"
	"fmt"
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
