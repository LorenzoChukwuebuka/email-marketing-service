package controller

import (
	"context"
	"email-marketing-service/core/handler/campaigns/dto"
	"email-marketing-service/core/handler/campaigns/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CampaignController struct {
	service *services.Service
}

func NewCampaignController(service *services.Service) *CampaignController {
	return &CampaignController{
		service: service,
	}
}

func (c *CampaignController) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.CampaignDTO

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	req.CompanyID = companyID
	req.UserId = userId

	result, err := c.service.CreateCampaign(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *CampaignController) GetAllCampaigns(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	page, pageSize, _, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	req := &dto.FetchCampaignDTO{
		UserID:    userId,
		CompanyID: companyID,
		Offset:    offset,
		Limit:     limit,
	}

	result, err := c.service.GetAllCampaigns(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *CampaignController) GetSingleCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	req := &dto.FetchCampaignDTO{
		UserID:     userId,
		CompanyID:  companyID,
		CampaignID: id,
	}

	result, err := c.service.GetSingleCampaign(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *CampaignController) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.CampaignDTO

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	req.UserId = userId
	req.CompanyID = companyID

	if err = c.service.UpdateCampaign(ctx, req, id); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "Campaign updated successfully")
}

func (c *CampaignController) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	req := &dto.FetchCampaignDTO{
		UserID:     userId,
		CompanyID:  companyID,
		CampaignID: id,
	}

	if err = c.service.DeleteCampaign(ctx, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "Campaign Deleted successfuly")
}

func (c *CampaignController) CreateCampaignGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.CampaignGroupDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	result, err := c.service.CreateCampaignGroup(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *CampaignController) SendCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	var req *dto.SendCampaignDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req.UserId = userId
	req.CompanyId = companyId

	result, err := c.service.SendCampaign(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}
