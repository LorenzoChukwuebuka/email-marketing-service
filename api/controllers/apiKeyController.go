package controllers

import (
	"email-marketing-service/api/services"
	"net/http"
	"github.com/golang-jwt/jwt"

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
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(float64)

	result, err := c.APIkeySVC.GenerateAPIKey(int(userId))

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, result)
}

func (c *ApiKeyController) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(float64)

	result, err := c.APIkeySVC.GetAPIKey(int(userId))

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *ApiKeyController) UpdateAPIKey(w http.ResponseWriter, r *http.Request) {

}
