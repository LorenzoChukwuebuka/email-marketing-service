package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"net/http"
	"github.com/golang-jwt/jwt"
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

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.UserId = userId

	result, err := c.TemplateSVC.CreateAndUpdateMarketingTemplate(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, result)

}

func (c *TemplateController) GetAllMarketingTemplates(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.TemplateSVC.GetAllMarketingTemplates(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, result)

}

func (c *TemplateController) GetAllTransactionalTemplates(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.TemplateSVC.GetAllTransactionalTemplates(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, result)
}
