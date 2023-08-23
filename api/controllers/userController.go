package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

//instantiate otp service

func NewUserService() *services.UserService {
	return &services.UserService{}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("jwtclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	username := claims["username"].(string)
	email := claims["email"].(string)

	response := fmt.Sprintf("Welcome, %s (%s)!", username, email)
	w.Write([]byte(response))

	fmt.Fprint(w, "Hello world")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.User

	utils.DecodeRequestBody(r, &reqdata)
	userCreateService, err := NewUserService().CreateUser(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, userCreateService)

}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.OTP

	utils.DecodeRequestBody(r, &reqdata)

	err := NewUserService().VerifyUser(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}
	utils.SuccessResponse(w, 200, "User has been successfully verifed")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.LoginModel

	utils.DecodeRequestBody(r, &reqdata)

	result, err := NewUserService().Login(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, result)
}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.ForgetPassword

	utils.DecodeRequestBody(r, &reqdata)

	err := NewUserService().ForgetPassword(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "email sent successfully")
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.ResetPassword

	utils.DecodeRequestBody(r, &reqdata)

	err := NewUserService().ResetPassword(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "password reset successfully")
}
