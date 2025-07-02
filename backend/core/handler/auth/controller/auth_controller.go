package controller

import (
	"context"
	"email-marketing-service/core/handler/auth/dto"
	"email-marketing-service/core/handler/auth/service"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"
)

type Controller struct {
	authSrv *service.Service
}

var (
	cfg = config.LoadEnv()
)

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

func (c *Controller) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)

	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	var req *dto.ChangePassword

	err = helper.DecodeRequestBody(r, &req)
	if err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	err = c.authSrv.ChangePassword(ctx, userId, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "password changed successfully")
}

// RefreshTokenHandler handles refreshing access tokens using the refresh token
func (c *Controller) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.RefreshAccessToken
	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	key := cfg.JWTKey

	// Parse and validate the refresh token
	claims, err := helper.ParseToken(reqdata.RefreshToken, []byte(key))
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("invalid refresh token"), nil)
		return
	}

	// Safe type assertions with validation
	userId, ok := claims["userId"].(string)
	if !ok {
		helper.ErrorResponse(w, fmt.Errorf("invalid userId in token"), nil)
		return
	}

	companyId, ok := claims["company_id"].(string)
	if !ok {
		helper.ErrorResponse(w, fmt.Errorf("invalid company_id in token"), nil)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		helper.ErrorResponse(w, fmt.Errorf("invalid username in token"), nil)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		helper.ErrorResponse(w, fmt.Errorf("invalid email in token"), nil)
		return
	}

	// Generate a new access token using the claims
	accessToken, err := helper.GenerateAccessToken(userId, companyId, username, email)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("failed to generate access token"), nil)
		return
	}

	// Send the new access token to the client
	res := map[string]string{
		"access_token": accessToken,
	}
	helper.SuccessResponse(w, 200, res)
}

func (c *Controller) GoogleLogin(http.ResponseWriter, *http.Request) {}

func (c *Controller) GoogleCallback(http.ResponseWriter, *http.Request) {}

func (channel *Controller) GoogleSignUp(http.ResponseWriter, *http.Request) {}

func (c *Controller) GoogleSignUpCallback(http.ResponseWriter, *http.Request) {}
