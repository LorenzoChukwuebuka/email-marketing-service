package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.User

	utils.DecodeRequestBody(r, &reqdata)

	userCreateService, err := services.CreateUser(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, userCreateService)

}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.OTP

	utils.DecodeRequestBody(r, &reqdata)

	err := services.VerifyUser(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "User has been successfully verifed")

}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {

}
