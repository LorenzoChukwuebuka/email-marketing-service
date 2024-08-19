package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
)

type CampaignController struct {
	CampaignSVC *services.CampaignService
}

func NewCampaignController(campaignSVC *services.CampaignService) *CampaignController {
	return &CampaignController{
		CampaignSVC: campaignSVC,
	}
}

func (c *CampaignController) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var reqdata dto.CampaignDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.UserId = userId

	result, err := c.CampaignSVC.CreateCampaign(&reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *CampaignController) GetAllCampaigns(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	result, err := c.CampaignSVC.GetAllCampaigns(userId, page, pageSize)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *CampaignController) GetSingleCampaign() {}

func (c *CampaignController) EditCampaign() {}

func (c *CampaignController) DeleteCampaign() {}
