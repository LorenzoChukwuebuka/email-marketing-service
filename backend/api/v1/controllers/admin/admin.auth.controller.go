package adminController

import (
	"email-marketing-service/api/v1/dto"
	adminservice "email-marketing-service/api/v1/services/admin"
	"email-marketing-service/api/v1/utils"
	"net/http"
)

var (
	config = utils.Config{}
	key    = config.JWTKey
)

type AdminController struct {
	AdminService *adminservice.AdminService
}

func NewAdminController(adminservice *adminservice.AdminService) *AdminController {
	return &AdminController{
		AdminService: adminservice,
	}
}

var (
	response = &utils.ApiResponse{}
)

func (c *AdminController) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.Admin

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.AdminService.CreateAdmin(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.AdminLogin

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.AdminService.AdminLogin(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *AdminController) ChangePassword(w http.ResponseWriter, r *http.Request) {

}

// RefreshTokenHandler handles refreshing access tokens using the refresh token
func (c *AdminController) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	var reqdata *dto.RefreshAccessToken

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	// Parse and validate the refresh token
	claims, err := utils.ParseToken(reqdata.RefreshToken, []byte(key))
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token using the claims (same user info)
	accessToken, err := utils.GenerateAdminAccessToken(claims["userId"].(string), claims["uuid"].(string), claims["type"].(string), claims["email"].(string))
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	// Send the new access token to the client
	res := map[string]string{
		"access_token": accessToken,
	}

	response.SuccessResponse(w, 200, res)
}
