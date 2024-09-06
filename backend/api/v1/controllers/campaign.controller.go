package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/mssola/user_agent"
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

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

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
	searchQuery := r.URL.Query().Get("search")

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

	result, err := c.CampaignSVC.GetAllCampaigns(userId, page, pageSize, searchQuery)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *CampaignController) GetAllScheduledCampaigns(w http.ResponseWriter, r *http.Request) {
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

	result, err := c.CampaignSVC.GetScheduledCampaigns(userId, page, pageSize)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *CampaignController) GetSingleCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.CampaignSVC.GetSingleCampaign(userId, campaignId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *CampaignController) EditCampaign(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.CampaignDTO
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	reqdata.UserId = userId
	reqdata.UUID = campaignId

	fmt.Printf("%+v\n",reqdata)

	if err := c.CampaignSVC.UpdateCampaign(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "campaign edited successfully")
}

func (c *CampaignController) AddOrEditCampaignGroup(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.CampaignGroupDTO
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	reqdata.UserId = userId

	if err := c.CampaignSVC.AddOrEditCampaignGroup(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "campaign group successfully")
}

func (c *CampaignController) SendCampaign(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.SendCampaignDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	reqdata.UserId = userId

	if err := c.CampaignSVC.SendCampaign(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "campaign sent successfully")

}

func (c *CampaignController) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	if err := c.CampaignSVC.DeleteCampaign(campaignId, userId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "campaign deleted successfully")
}

func (c *CampaignController) TrackOpenCampaignEmails(w http.ResponseWriter, r *http.Request) {
	contactMail := r.URL.Query().Get("email")
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]

	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	ua := user_agent.New(r.Header.Get("User-Agent"))

	deviceType := "desktop"

	if ua.Mobile() {
		deviceType = "mobile"
	}

	osInfo := ua.OSInfo()
	if osInfo.Name == "iPad" || osInfo.Name == "Android" {
		deviceType = "tablet"
	}

	if err := c.CampaignSVC.TrackOpenCampaignEmails(campaignId, contactMail, deviceType, ipAddress); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "campaign tracked successfully")
}

func (c *CampaignController) TrackClickedCampaignsEmails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]
	contactMail := r.URL.Query().Get("email")
	if err := c.CampaignSVC.TrackClickedCampaignsEmails(campaignId, contactMail); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "campaign tracked successfully")
}

func (c *CampaignController) UnsubscribeFromCampaign(w http.ResponseWriter, r *http.Request) {

	contactMail := r.URL.Query().Get("email")
	camapignId := r.URL.Query().Get("campaign")

	if err := c.CampaignSVC.UnsubscribeFromCampaign(camapignId, contactMail); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	templatePath := filepath.Join("api", "v1", "templates", "unsubscribe.templ")
	unsubscribeTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(unsubscribeTemplate)
}

func (c *CampaignController) GetEmailResultStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]

	result, err := c.CampaignSVC.GetEmailResultStats(campaignId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *CampaignController) GetAllRecipientsForACampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId := vars["campaignId"]

	result, err := c.CampaignSVC.GetAllRecipientsForACampaign(campaignId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *CampaignController) GetUserCampaignStats(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.CampaignSVC.GetUserCampaignStats(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *CampaignController) GetUserCampaignsStats(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.CampaignSVC.GetUserCampaignsStats(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}
