package controllers

import (
	"email-marketing-service/api/services"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type ApiKeyController struct {
	APIkeySVC *services.ApiKeyService
}

func NewAPIKeyController(apiKeyService *services.ApiKeyService) *ApiKeyController {
	return &ApiKeyController{}
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
