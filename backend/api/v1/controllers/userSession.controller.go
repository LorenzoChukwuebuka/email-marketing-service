package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"net/http"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
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
	var reqdata *dto.UserSession

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.UserSessionSVC.CreateSession(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserSessionController) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.UserSessionSVC.GetAllSessions(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserSessionController) DeleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sessionId := vars["subscriptionId"]

	if err := c.UserSessionSVC.DeleteSession(sessionId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "session deleted successfully")

}
