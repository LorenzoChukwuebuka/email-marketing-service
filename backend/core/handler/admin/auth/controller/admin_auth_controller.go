package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/auth/dto"
	"email-marketing-service/core/handler/admin/auth/services"
	"email-marketing-service/internal/config"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"
)

var (
	key = config.LoadEnv().JWTKey
)

type Controller struct {
	adminauthsvc *services.Service
}

func NewAdminAuthController(adminauthsvc *services.Service) *Controller {
	return &Controller{
		adminauthsvc: adminauthsvc,
	}
}

func (c *Controller) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.AdminRequestDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	result, err := c.adminauthsvc.CreateAdmin(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 201, result)
}

func (c *Controller) AdminLogin(w http.ResponseWriter, r *http.Request) {

}

// RefreshTokenHandler handles refreshing access tokens using the refresh token
func (c *Controller) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	var reqdata *dto.RefreshAccessToken

	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	// Parse and validate the refresh token
	claims, err := helper.ParseToken(reqdata.RefreshToken, []byte(key))
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token using the claims (same user info)
	accessToken, err := helper.GenerateAdminAccessToken(claims["userId"].(string), claims["uuid"].(string), claims["type"].(string), claims["email"].(string))
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	// Send the new access token to the client
	res := map[string]string{
		"access_token": accessToken,
	}

	helper.SuccessResponse(w, 200, res)
}
