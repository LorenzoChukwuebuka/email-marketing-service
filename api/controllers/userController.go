package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Welcome(w http.ResponseWriter, r *http.Request) {
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

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.User

	utils.DecodeRequestBody(r, &reqdata)
	userCreateService, err := c.userService.CreateUser(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, userCreateService)

}

func (c *UserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.OTP

	utils.DecodeRequestBody(r, &reqdata)

	err := c.userService.VerifyUser(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}
	utils.SuccessResponse(w, 200, "User has been successfully verifed")
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.LoginModel

	utils.DecodeRequestBody(r, &reqdata)

	result, err := c.userService.Login(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, result)
}

func (c *UserController) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.ForgetPassword

	utils.DecodeRequestBody(r, &reqdata)

	err := c.userService.ForgetPassword(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "email sent successfully")
}

func (c *UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.ResetPassword

	utils.DecodeRequestBody(r, &reqdata)

	err := c.userService.ResetPassword(reqdata)

	if err != nil {
		utils.ErrorResponse(w, err.Error())
		return
	}

	utils.SuccessResponse(w, 200, "password reset successfully")
}
