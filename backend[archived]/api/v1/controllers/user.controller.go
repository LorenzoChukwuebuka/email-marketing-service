package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"email-marketing-service/internals/workers/payloads"
	"email-marketing-service/internals/workers/tasks"
	"email-marketing-service/pkg/asynqpkg"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/hibiken/asynq"
	"log"
	"net/http"
	"os"
	"time"
)

type UserController struct {
	userService *services.UserService
	loginAuth   *utils.GoogleAuthConfig
	signupAuth  *utils.GoogleAuthConfig
	client      *asynq.Client
}

func NewUserController(userService *services.UserService) *UserController {
	// Determine which redirect URL to use based on environment
	loginRedirectURL := config.GOOGLE_CLIENT_LOGIN_REDIRECT_URL
	signupRedirectURL := config.GOOGLE_CLIENT_SIGNUP_REDIRECT_URL

	if os.Getenv("SERVER_MODE") == "production" {
		loginRedirectURL = config.GOOGLE_CLIENT_LOGIN_REDIRECT_URL_PROD
		signupRedirectURL = config.GOOGLE_CLIENT_SIGNUP_REDIRECT_URL_PROD
	}

	loginAuth := utils.NewGoogleAuth(
		utils.WithClientID(config.GOOGLE_CLIENT_ID),
		utils.WithClientSecret(config.GOOGLE_CLIENT_SECRET),
		utils.WithRedirectURL(loginRedirectURL),
	)

	signupAuth := utils.NewGoogleAuth(
		utils.WithClientID(config.GOOGLE_CLIENT_ID),
		utils.WithClientSecret(config.GOOGLE_CLIENT_SECRET),
		utils.WithRedirectURL(signupRedirectURL),
	)

	client := asynqpkg.GetClient()

	return &UserController{
		userService: userService,
		loginAuth:   loginAuth,
		signupAuth:  signupAuth,
		client:      client,
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

func (c *UserController) GetClient() *asynq.Client {
	return c.client
}

func (c *UserController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()
	authURL := c.loginAuth.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (c *UserController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	userData, err := c.loginAuth.GetUserData(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	redirectURL := fmt.Sprintf("%sauth/callback?token=%s&type=%s", config.FRONTEND_URL, jwtToken, "login")
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (c *UserController) GoogleSignup(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()
	authURL := c.signupAuth.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (c *UserController) GoogleSignUpCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	userData, err := c.signupAuth.GetUserData(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDto := &dto.User{
		FullName: userData.Name,
		Email:    userData.Email,
		GoogleID: userData.ID,
	}

	// First try to login if user exists
	loginData, err := c.userService.GoogleLogin(userDto)
	if err == nil {
		//save the user details in the asynq store this will come in handy when we get to the FE

		// Store user details in Asynq with expiration
		userResponse := loginData["details"].(model.UserResponse)
		task, err := tasks.NewStoreUserDetailsTask(utils.ToMap(userResponse))
		if err != nil {
			http.Error(w, "Error creating task", http.StatusInternalServerError)
			return
		}

		// Store with 5 minute expiration
		info, err := c.client.Enqueue(task,
			asynq.Queue("user_details"),
			asynq.Timeout(5*time.Minute),
			asynq.Retention(5*time.Minute),
		)

		if err != nil {
			http.Error(w, "Error storing details", http.StatusInternalServerError)
			return
		}

		// User exists, return login data
		redirectURL := fmt.Sprintf("%sauth/callback?type=login&token=%s&refresh_token=%s&details_key=%s",
			config.FRONTEND_URL,
			loginData["token"],
			loginData["refresh_token"],
			info.ID) // Using task ID as the key
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	// If user doesn't exist, create new user and return login data
	signupData, err := c.userService.CreateGoogleUser(userDto)
	if err != nil {
		http.Error(w, "error signing up: "+err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("%sauth/callback?type=signup&token=%s&refresh_token=%s",
		config.FRONTEND_URL,
		signupData["token"],
		signupData["refresh_token"])
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (c *UserController) GetLoginDetails(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		response.ErrorResponse(w, "No details key provided")
		return
	}

	// Create inspector using the same Redis connection as your client
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: config.REDIS_PORT,
	})

	taskInfo, err := inspector.GetTaskInfo("user_details", key)
	if err != nil {
		response.ErrorResponse(w, "Details not found or expired")
		return
	}

	var payload payloads.UserDetailsPayload
	if err := json.Unmarshal(taskInfo.Payload, &payload); err != nil {
		response.ErrorResponse(w, "Error decoding details")
		return
	}

	// Delete the task after fetching
	if err := inspector.DeleteTask("user_details", key); err != nil {
		log.Printf("Error deleting task: %v", err)
	}

	response.SuccessResponse(w, http.StatusOK, payload.Details)
}
