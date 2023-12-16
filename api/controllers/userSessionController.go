package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"net/http"
)

type UserSessionController struct {
	UserSessionSVC *services.UserSessionService
}

func NewUserSessionController(usersessionSvc *services.UserSessionService) *UserSessionController {
	return &UserSessionController{
		UserSessionSVC: usersessionSvc,
	}
}

func (c *UserSessionController) CreateSessions(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.UserSessionModelStruct

	utils.DecodeRequestBody(r, &reqdata)
 
	result,err := c.UserSessionSVC.CreateSession(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return 
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserSessionController) GetAllSessions(w http.ResponseWriter, r *http.Request) {

}

func (c *UserSessionController) DeleteSession(w http.Response, r *http.Request) {

}
