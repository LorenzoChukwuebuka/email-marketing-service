package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type SMTPKeyController struct {
	SMTPkeyService *services.SMTPKeyService
}

func NewSMTPKeyController(smtpkeyService *services.SMTPKeyService) *SMTPKeyController {
	return &SMTPKeyController{
		SMTPkeyService: smtpkeyService,
	}
}

func (c *SMTPKeyController) GenerateNewSMTPMasterPassword(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	result, err := c.SMTPkeyService.GenerateNewSMTPMasterPassword(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *SMTPKeyController) GetUserSMTPKeys(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	result, err := c.SMTPkeyService.GetSMTPKeys(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *SMTPKeyController) CreateSMTPKey(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	var reqdata *dto.SMTPKeyDTO

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	reqdata.UserId = userId
	result, err := c.SMTPkeyService.CreateSMTPKey(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *SMTPKeyController) ToggleSMTPKeyStatus(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	vars := mux.Vars(r)
	smtpkeyId := vars["smtpKeyId"]
	if err := c.SMTPkeyService.ToggleSMTPKeyStatus(userId, smtpkeyId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "toggle successful")
}

func (c *SMTPKeyController) DeleteSMTPKey(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	smtpkeyId := vars["smtpKeyId"]

	if err := c.SMTPkeyService.DeleteSMTPKey(smtpkeyId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "smtp key deleted successfully")
}
