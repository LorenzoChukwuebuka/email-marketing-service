package controller

import (
	"context"
	"email-marketing-service/core/handler/api_smtp_keys/dto"
	"email-marketing-service/core/handler/api_smtp_keys/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type SMTPKeyController struct {
	smtpKeyService *services.SMTPKeyService
}

func NewSMTPKeyController(smtpKeyService *services.SMTPKeyService) *SMTPKeyController {
	return &SMTPKeyController{
		smtpKeyService: smtpKeyService,
	}
}

func (c *SMTPKeyController) CreateSMTPKEY(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.SMTPKeyRequestDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	req.CompanyID = companyId
	req.UserId = userId

	result, err := c.smtpKeyService.CreateSMTPKey(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 201, result)

}

func (c *SMTPKeyController) GenerateNewSMTPMasterPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	result, err := c.smtpKeyService.GenerateNewSMTPMasterPassword(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 201, result)
}

func (c *SMTPKeyController) GetSMTPKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	result, err := c.smtpKeyService.GetSMTPKeys(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *SMTPKeyController) ToggleSMTPKeyStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	vars := mux.Vars(r)
	smtpKeyId := vars["smtpkeyId"]

	if err := c.smtpKeyService.ToggleSMTPKeyStatus(ctx, userId, smtpKeyId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "smtp key status updated successfully")
}

func (c *SMTPKeyController) DeleteSMTPKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	smtpKeyId := vars["smtpkeyId"]

	if err := c.smtpKeyService.DeleteSMTPKey(ctx, smtpKeyId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "smtp key deleted successfully")
}
