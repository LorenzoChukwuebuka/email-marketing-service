package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"net/http"
)

type SMTPMailController struct {
	APISVC *services.APIKeyService
}

func NewSMTPMailController(apiservice *services.APIKeyService) *SMTPMailController {
	return &SMTPMailController{
		APISVC: apiservice,
	}
}

func (c *SMTPMailController) SendSMTPMail(w http.ResponseWriter, r *http.Request) {

	var reqdata *model.EmailRequest

	utils.DecodeRequestBody(r, &reqdata)

	// Get the value of the "api-key" header
	apiKey := r.Header.Get("api-key")
	if apiKey == "" {
		// The header is not present or has an empty value
		errorRes := map[string]interface{}{
			"status":         http.StatusUnauthorized,
			"error response": "API key not provided",
		}
		response.ErrorResponse(w, errorRes)
		return
	}

	apiKeyExist, err := c.APISVC.APIKeyRepo.CheckIfAPIKEYExists(apiKey)

	if err != nil {
		errorres := map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "error fetching api key",
		}

		response.ErrorResponse(w, errorres)

		return
	}

	if !apiKeyExist {
		errorRes := map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"message": "Invalid API key provided",
		}
		response.ErrorResponse(w, errorRes)
		return
	}

	response.SuccessResponse(w, 200, reqdata)

}
