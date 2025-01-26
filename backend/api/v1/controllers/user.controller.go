package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	config   = utils.LoadEnv()
	key      = config.JWTKey
	response = &utils.ApiResponse{}
)

// Define at package level
var loginOauthConfig *oauth2.Config
var signupOauthConfig *oauth2.Config

// Initialize in init() function
func init() {
	loginOauthConfig = &oauth2.Config{
		ClientID:     config.GOOGLE_CLIENT_ID,
		ClientSecret: config.GOOGLE_CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_CLIENT_LOGIN_REDIRECT_URL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	signupOauthConfig = &oauth2.Config{
		ClientID:     config.GOOGLE_CLIENT_ID,
		ClientSecret: config.GOOGLE_CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_CLIENT_SIGNUP_REDIRECT_URL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

}

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

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

func (c *UserController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":  200,
		"message": "The application is working ....",
	}
	response.SuccessResponse(w, 200, health)
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.User

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	userCreateService, err := c.userService.CreateUser(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userCreateService)
}

func (c *UserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.OTP

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	err := c.userService.VerifyUser(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "User has been successfully verifed")
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.Login

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	result, err := c.userService.Login(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserController) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.ResendOTP

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	if err := c.userService.ResendOTP(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "otp resent successfully")
}

func (c *UserController) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.ForgetPassword

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	if err := c.userService.ForgetPassword(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "email sent successfully")
}

func (c *UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.ResetPassword

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	if err := c.userService.ResetPassword(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "password reset successfully")
}

func (c *UserController) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	var reqdata *dto.ChangePassword
	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	err = c.userService.ChangePassword(userId, reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "password changed successfully")
}

func (c *UserController) EditUser(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	var reqdata *model.User

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	reqdata.UUID = userId

	if err := c.userService.EditUser(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "User edited successfully")
}

func (c *UserController) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	result, err := c.userService.GetUserDetails(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserController) GetUserSubscription(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	result, err := c.userService.GetUserCurrentRunningSubscriptionWithMailsRemaining(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *UserController) GetAllUserNotifications(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	result, err := c.userService.GetAllNotifications(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

// this is just for testing purposes....
func (c *UserController) SSEGetAllUserNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userId, err := ExtractUserId(r)
	if err != nil {
		WriteSSEError(w, err)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		WriteSSEError(w, fmt.Errorf("streaming not supported"))
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			notifications, err := c.userService.GetAllNotifications(userId)
			if err != nil {
				WriteSSEError(w, err)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", ToJSON(notifications))
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

func (c *UserController) UpdateReadStatus(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	err = c.userService.UpdateReadStatus(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "notifications updated")
}

func (c *UserController) MarkUserForDeletion(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	if err := c.userService.MarkUserForDeletion(userId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

}

func (c *UserController) CancelUserDeletion(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}

	if err := c.userService.CancelUserDeletion(userId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
}

// RefreshTokenHandler handles refreshing access tokens using the refresh token
func (c *UserController) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.RefreshAccessToken
	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}

	// Parse and validate the refresh token
	claims, err := utils.ParseToken(reqdata.RefreshToken, []byte(key))
	if err != nil {
		response.ErrorResponse(w, "Invalid refresh token")
		return
	}
	// Generate a new access token using the claims (same user info)
	accessToken, err := utils.GenerateAccessToken(claims["userId"].(string), claims["uuid"].(string), claims["username"].(string), claims["email"].(string))
	if err != nil {
		response.ErrorResponse(w, "Failed to generate access token")
		return
	}
	// Send the new access token to the client
	res := map[string]string{
		"access_token": accessToken,
	}
	response.SuccessResponse(w, 200, res)
}

func (c *UserController) GoogleSignup(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()
	// Generate the Google OAuth2 URL and redirect to it
	authURL := loginOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (c *UserController) GoogleSignUpCallback(w http.ResponseWriter, r *http.Request) {
	// Get the code from the callback
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	token, err := signupOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info
	client := signupOauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData := struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		VerifiedEmail bool   `json:"verified_email"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to decode user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userDto := &dto.User{
		FullName: userData.Name,
		Email:    userData.Email,
		GoogleID: userData.ID,
	}

	_, err = c.userService.CreateUser(userDto)

	if err != nil {
		log.Printf("User creation error: %v", err)
		http.Error(w, "error signing up", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("%sauth/signup/callback", config.FRONTEND_URL)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (c *UserController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate the Google OAuth2 URL and redirect to it
	state := utils.GenerateRandomState()
	authURL := loginOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the Google OAuth callback
func (c *UserController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Get the code from the callback
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	token, err := loginOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info
	client := loginOauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData := struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		VerifiedEmail bool   `json:"verified_email"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to decode user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	print(userData.ID)

	// Generate JWT token
	jwtToken, err := utils.GenerateAccessToken(
		userData.ID,
		"google",
		userData.Name,
		userData.Email,
	)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to frontend with the token
	redirectURL := fmt.Sprintf("%sauth/callback?token=%s", config.FRONTEND_URL, jwtToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
