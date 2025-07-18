package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/campaigns/dto"
	"email-marketing-service/core/handler/admin/campaigns/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type AdminCampaignController struct {
	adminCampaignServices *services.AdminCampaignService
}

func NewAdminCampaignController(adminCampaignServices *services.AdminCampaignService) *AdminCampaignController {
	return &AdminCampaignController{
		adminCampaignServices: adminCampaignServices,
	}
}

func (c *AdminCampaignController) GetAllUserCampaigns(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]
	companyId := vars["companyId"]

	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	req := &dto.AdminFetchCampaignDTO{
		Search:    search,
		Offset:    offset,
		Limit:     limit,
		UserID:    userId,
		CompanyID: companyId,
	}

	result, err := c.adminCampaignServices.GetAllUserCampaigns(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *AdminCampaignController) GetSingleCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["campaignId"]
	userId := vars["userId"]
	companyID := vars["companyId"]

	req := &dto.AdminFetchCampaignDTO{
		UserID:     userId,
		CompanyID:  companyID,
		CampaignID: id,
	}

	result, err := c.adminCampaignServices.GetSingleCampaign(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminCampaignController) GetAllRecipientsForACampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()	

	vars := mux.Vars(r)
	campaignId := vars["campaignId"]
	companyId := vars["companyId"]

	result, err := c.adminCampaignServices.GetAllRecipientsForACampaign(ctx, campaignId, companyId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}
