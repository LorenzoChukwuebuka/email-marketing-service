package adminController

import (
	"email-marketing-service/api/v1/repository"
	adminservice "email-marketing-service/api/v1/services/admin"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AdminCampaignController struct {
	AdminCampaignSVC *adminservice.AdminCampaignService
}

func NewAdminCampaginController(adminCampaignSVC *adminservice.AdminCampaignService) *AdminCampaignController {
	return &AdminCampaignController{
		AdminCampaignSVC: adminCampaignSVC,
	}
}

func (c *AdminCampaignController) GetAllUserCampaigns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	searchQuery := r.URL.Query().Get("search")

	params := repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	campaigns, err := c.AdminCampaignSVC.GetAllUserCampaigns(userId, searchQuery, params)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, campaigns)
}

func (c *AdminCampaignController) GetAUserCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignUUID := vars["campaignUUID"]

	campaign, err := c.AdminCampaignSVC.GetAUserCampaign(campaignUUID)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, campaign)
}

func (c *AdminCampaignController) DeleteUserCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignUUID := vars["campaignUUID"]

	err := c.AdminCampaignSVC.DeleteUserCampaign(campaignUUID)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Campaign successfully deleted")
}

func (c *AdminCampaignController) SuspendUserCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignUUID := vars["campaignUUID"]

	err := c.AdminCampaignSVC.SuspendUserCampaign(campaignUUID)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Campaign successfully suspended")
}
