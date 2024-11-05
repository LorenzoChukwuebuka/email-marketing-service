package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/gorilla/mux"
	"net/http"
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
	var reqdata *dto.APIkeyDTO
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	reqdata.UserId = userId
	result, err := c.APIkeySVC.GenerateAPIKey(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 201, result)
}

func (c *ApiKeyController) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	result, err := c.APIkeySVC.GetAPIKey(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *ApiKeyController) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apikeyId := vars["apiKeyId"]
	err := c.APIkeySVC.DeleteAPIKey(apikeyId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "api key deleted successfully")
}
