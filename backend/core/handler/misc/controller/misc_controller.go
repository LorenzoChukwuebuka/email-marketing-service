package controller

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/misc/services"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mssola/user_agent"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type MiscController struct {
	miscservice *services.MiscService
	store       db.Store
}

func NewMiscController(miscservice *services.MiscService, store db.Store) *MiscController {
	return &MiscController{
		miscservice: miscservice,
		store:       store,
	}
}

func (c *MiscController) TrackOpenCampaignEmails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

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

	if err := c.miscservice.TrackOpenCampaignEmails(ctx, campaignId, contactMail, deviceType, ipAddress); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "campaign tracked successfully")
}

func (c *MiscController) TrackClickedCampaignsEmails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	campaignId := vars["campaignId"]
	contactMail := r.URL.Query().Get("email")

	if err := c.miscservice.TrackClickedCampaignsEmails(ctx, campaignId, contactMail); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "campaign tracked successfully")
}

func (c *MiscController) UnsubscribeFromCampaign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	contactMail := r.URL.Query().Get("email")
	camapignId := r.URL.Query().Get("campaign")
	companyId := r.URL.Query().Get("companyId")

	if err := c.miscservice.UnsubscribeFromCampaign(ctx, camapignId, contactMail, companyId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	templatePath := filepath.Join("internal", "templates", "unsubscribe.templ")
	unsubscribeTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(unsubscribeTemplate)
}

func (c *MiscController) SendSMTPMail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var reqdata *domain.EmailRequest

	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	// Get the value of the "api-key" header
	apiKey := r.Header.Get("api-key")
	if apiKey == "" {
		// The header is not present or has an empty value
		helper.ErrorResponse(w, fmt.Errorf("API key not provided"), http.StatusUnauthorized)
		return
	}

	key, err := c.store.CheckIfAPIKeyExists(ctx, apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.ErrorResponse(w, fmt.Errorf("invalid API Key"), nil)
			return
		}
		helper.ErrorResponse(w, common.ErrFetchingRecord, nil)
		return
	}

	result, err := c.miscservice.PrepareMail(ctx, reqdata, key.CompanyID, key.UserID)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *MiscController) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	plans, err := c.miscservice.GetAllPlansWithDetails(ctx)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, plans)
}

func (c *MiscController) GetPlanByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)

	planId := vars["planId"]

	planID, err := common.ParseUUIDMap(map[string]string{"planId": planId})
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	plan, err := c.miscservice.GetSinglePlan(ctx, planID["planId"])
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, plan)
}
