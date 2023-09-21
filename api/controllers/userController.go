package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

var (
	response = &utils.ApiResponse{}
)

func (c *UserController) Welcome(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
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

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.User
	
	utils.DecodeRequestBody(r, &reqdata)
	userCreateService, err := c.userService.CreateUser(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userCreateService)

}

func (c *UserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.OTP

	utils.DecodeRequestBody(r, &reqdata)

	err := c.userService.VerifyUser(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "User has been successfully verifed")
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.LoginModel

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.userService.Login(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserController) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.ForgetPassword

	utils.DecodeRequestBody(r, &reqdata)

	if err := c.userService.ForgetPassword(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "email sent successfully")
}

func (c *UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.ResetPassword

	utils.DecodeRequestBody(r, &reqdata)

	if err := c.userService.ResetPassword(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "password reset successfully")
}

func (c *UserController) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(float64)

	var reqdata *model.ChangePassword

	utils.DecodeRequestBody(r, &reqdata)

	err := c.userService.ChangePassword(int(userId), reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "password changed successfully")

}

func (c *UserController) EditUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(float64)

	var reqdata *model.User

	utils.DecodeRequestBody(r, &reqdata)

	if err := c.userService.EditUser(int(userId),reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "User edited successfully")

}
