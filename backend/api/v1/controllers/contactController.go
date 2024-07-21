package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type ContactController struct {
	ContactService *services.ContactService
}

func NewContactController() *ContactController {
	return &ContactController{}
}

func (c *ContactController) CreateContact(w http.ResponseWriter, r *http.Request) {
	var reqdata dto.ContactDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.UserId = userId

	result, err := c.ContactService.CreateContact(&reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}
