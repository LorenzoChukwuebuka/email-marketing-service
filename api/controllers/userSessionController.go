package controllers

import (
	"net/http"
)

type UserSessionController struct{}

func NewUserSessionController() *UserSessionController {
	return &UserSessionController{}
}

func (c *UserSessionController) GetAllSessions(w http.Response, r *http.Request) {

}

func (c *UserSessionController) DeleteSession(w http.Response, r *http.Request) {

}
