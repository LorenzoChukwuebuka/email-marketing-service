package adminController

import (
	adminmodel "email-marketing-service/api/model/admin"
	adminservice "email-marketing-service/api/services/admin"
	"email-marketing-service/api/utils"
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
	var reqdata *adminmodel.AdminModel

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.AdminService.CreateAdmin(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	var reqdata *adminmodel.AdminLogin

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.AdminService.AdminLogin(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}
