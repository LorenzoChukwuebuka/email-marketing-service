package controller

import (
	"context"
	"email-marketing-service/core/handler/auth/dto"
	"email-marketing-service/core/handler/auth/service"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"
)

type Controller struct {
	authSrv *service.Service
}

func NewAuthController(authSrv *service.Service) *Controller {
	return &Controller{
		authSrv: authSrv,
	}
}

func (c *Controller) Welcome(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("Welcome, %s (%s)!", "hello@hello.com", "hello@hello.com")
	w.Write([]byte(response))
	fmt.Fprint(w, "Hello world")
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.LoginRequest

	err := helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	result, err := c.authSrv.LoginUser(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 201, result)
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.UserSignUpRequest

	err := helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	result, err := c.authSrv.SignUp(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 201, result)
}

func (c *Controller) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.VerifyUserRequest

	err := helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	err = c.authSrv.VerifyUser(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "email verified successfully")
}

func (c *Controller) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.ForgetPassword

	err := helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	err = c.authSrv.ForgetPassword(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 201, "email sent successfully")
}

func (c *Controller) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.ResetPassword

	err := helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	err = c.authSrv.ResetPassword(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "password reset successful")
}

func (c *Controller) GoogleLogin(http.ResponseWriter, *http.Request) {}

func (c *Controller) GoogleCallback(http.ResponseWriter, *http.Request) {}

func (channel *Controller) GoogleSignUp(http.ResponseWriter, *http.Request) {}

func (c *Controller) GoogleSignUpCallback(http.ResponseWriter, *http.Request) {}
