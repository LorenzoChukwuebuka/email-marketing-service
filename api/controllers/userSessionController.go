package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
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

func (c *UserSessionController) getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func (c *UserSessionController) CreateSessions(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.UserSession

	ipAddress := c.getIPAddress(r)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.IPAddress = &ipAddress

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

	userId := claims["userId"].(float64)

	result, err := c.UserSessionSVC.GetAllSessions(int(userId))

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
