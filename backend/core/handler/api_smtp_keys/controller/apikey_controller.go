package controller

import (
	"context"
	"email-marketing-service/core/handler/api_smtp_keys/dto"
	"email-marketing-service/core/handler/api_smtp_keys/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ApiKeyController struct {
	APIkeySVC *services.APIKeyService
}

func NewAPIKeyController(apiKeyService *services.APIKeyService) *ApiKeyController {
	return &ApiKeyController{
		APIkeySVC: apiKeyService,
	}
}

func (c *ApiKeyController) GenerateAPIKEY(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.APIkeyRequestDTO
	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	err = helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	req.UserId = userId
	req.CompanyID = companyId
	result, err := c.APIkeySVC.GenerateAPIKey(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 201, result)
}

func (c *ApiKeyController) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	result, err := c.APIkeySVC.GetAPIKey(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *ApiKeyController) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	apikeyId := vars["apikeyId"]

	key, err := uuid.Parse(apikeyId)
	if err != nil {
		helper.ErrorResponse(w, common.ErrInvalidUUID, nil)
		return
	}

	err = c.APIkeySVC.DeleteAPIKey(ctx, key)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "api key deleted successfully")
}
