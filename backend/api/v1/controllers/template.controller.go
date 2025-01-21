package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type TemplateController struct {
	TemplateSVC services.TemplateService
}

func NewTemplateController(templateSvc *services.TemplateService) *TemplateController {
	return &TemplateController{
		TemplateSVC: *templateSvc,
	}
}

func (c *TemplateController) CreateAndUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.TemplateDTO

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
	result, err := c.TemplateSVC.CreateTemplate(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *TemplateController) GetAllMarketingTemplates(w http.ResponseWriter, r *http.Request) {
	page, pageSize, searchQuery, err := ParsePaginationParams(r)

	if err != nil {
		HandleControllerError(w, err)
	}

	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	result, err := c.TemplateSVC.GetAllMarketingTemplates(userId, page, pageSize, searchQuery)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *TemplateController) GetAllTransactionalTemplates(w http.ResponseWriter, r *http.Request) {

	page, pageSize, searchQuery, err := ParsePaginationParams(r)

	if err != nil {
		HandleControllerError(w, err)
	}

	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	result, err := c.TemplateSVC.GetAllTransactionalTemplates(userId, page, pageSize, searchQuery)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *TemplateController) GetTransactionalTemplate(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	vars := mux.Vars(r)

	templateId := vars["templateId"]

	result, err := c.TemplateSVC.GetTransactionalTemplate(userId, templateId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *TemplateController) GetMarketingTemplate(w http.ResponseWriter, r *http.Request) {

	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	vars := mux.Vars(r)

	templateId := vars["templateId"]

	result, err := c.TemplateSVC.GetMarketingTemplate(userId, templateId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *TemplateController) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.TemplateDTO
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	vars := mux.Vars(r)

	templateId := vars["templateId"]

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	reqdata.UserId = userId

	if err := c.TemplateSVC.UpdateTemplate(reqdata, templateId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "template updated successfully")
}

func (c *TemplateController) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	vars := mux.Vars(r)

	templateId := vars["templateId"]

	if err := c.TemplateSVC.DeleteTemplate(userId, templateId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "template deleted successfully")
}

func (c *TemplateController) SendTestMail(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.SendTestMailDTO
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

	if err := c.TemplateSVC.SendTestMail(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "mails sent successfully")

}
