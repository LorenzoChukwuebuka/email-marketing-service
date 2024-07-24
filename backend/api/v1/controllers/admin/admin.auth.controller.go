package adminController

import (
	"email-marketing-service/api/v1/dto"
	adminservice "email-marketing-service/api/v1/services/admin"
	"email-marketing-service/api/v1/utils"
	"net/http"
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
