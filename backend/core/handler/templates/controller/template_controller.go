package controller

import (
	"context"
	"email-marketing-service/core/handler/templates/dto"
	"email-marketing-service/core/handler/templates/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Controller struct {
	service *services.Service
}

func NewTemplateController(service *services.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	var req dto.TemplateDTO
	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	req.UserId = userId
	req.CompanyID = companyId

	template, err := c.service.CreateTemplate(ctx, &req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, 201, template)
}

func (c *Controller) GetTemplatesByType(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	tempType := vars["type"]

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	page, pageSize, _, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	if err != nil {
		helper.ErrorResponse(w, err, "an error occured")
		return
	}

	req := &dto.FetchTemplateDTO{
		UserId: userId,
		Type:   tempType,
		Offset: offset,
		Limit:  limit,
	}

	template, err := c.service.GetAllTemplateByType(ctx, *req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, 200, template)
}

func (c *Controller) GetTemplateById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	tempId := vars["id"]
	temptype := vars["type"]

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req := &dto.FetchTemplateDTO{
		UserId:     userId,
		TemplateId: tempId,
		Type:       temptype,
	}

	fmt.Printf("ddd %+v", req)
	template, err := c.service.GetTemplateByID(ctx, *req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, 200, template)
}

func (c *Controller) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	tempId := vars["id"]

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	var req dto.TemplateDTO
	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	req.UserId = userId
	req.TemplateID = tempId
	req.CompanyID = companyId

	template, err := c.service.UpdateTemplate(ctx, &req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, 200, template)
}

func (c *Controller) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	tempId := vars["id"]

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req := &dto.FetchTemplateDTO{
		UserId:     userId,
		TemplateId: tempId,
	}

	err = c.service.DeleteTemplate(ctx, *req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, 200, "successfully deleted template")
}

func (c *Controller) SendTestMail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	var req dto.SendTestMailDTO
	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	req.UserId = userId
	err = c.service.SendTestMail(ctx, &req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, 200, "Test mail sent successfully")
}
