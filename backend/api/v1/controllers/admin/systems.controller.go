package adminController

import (
	"email-marketing-service/api/v1/dto"
	adminservice "email-marketing-service/api/v1/services/admin"
	"email-marketing-service/api/v1/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type SystemsController struct {
	SystemsService *adminservice.SystemsService
}

func NewSystemsController(systemsSvc *adminservice.SystemsService) *SystemsController {
	return &SystemsController{
		SystemsService: systemsSvc,
	}
}

func (c *SystemsController) CreateRecords(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.SystemsDTO

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	result, err := c.SystemsService.GenerateAndSaveSMTPCredentials(reqdata.Domain)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, result)
}

func (c *SystemsController) GetDNSRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	domainName := vars["domain"]

	result, err :=c.SystemsService.GetDNSRecords(domainName)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *SystemsController) DeleteDNSRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	domainName := vars["domain"]

	err := c.SystemsService.DeleteDNSRecords(domainName)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "deleted successfully")
}
