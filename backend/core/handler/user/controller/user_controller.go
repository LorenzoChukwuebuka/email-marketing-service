package controller

import (
	"context"
	"email-marketing-service/core/handler/auth/dto"
	"email-marketing-service/core/handler/user/services"

	//"email-marketing-service/internal/common"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	result, err := c.userService.GetUserNotifications(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *UserController) GetUserNotificationsLongPoll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	// Optional: Get 'sinceId' from query params
	// Client sends the UUID of their last received notification
	var sinceID *string
	if sinceIDStr := r.URL.Query().Get("sinceId"); sinceIDStr != "" {
		sinceID = &sinceIDStr
	}

	// Set headers for long polling
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	result, err := c.userService.GetUserNotificationsLongPoll(ctx, userId, sinceID)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *UserController) UpdateReadStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	err = c.userService.UpdateReadStatus(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "Notification marked as read")
}

func (c *UserController) MarkUserForDeletion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	err = c.userService.MarkUserForDeletion(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "User marked for deletion")
}

func (c *UserController) CancelUserDeletion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	err = c.userService.CancelUserDeletion(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "User deletion cancelled")
}

func (c *UserController) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	result, err := c.userService.GetUserDetails(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *UserController) GetCurrentRunningSubscription(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	_, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch company id from jwt"), nil)
		return
	}

	result, err := c.userService.GetCurrentRunningSubscription(ctx, companyId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *UserController) EditUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.EditUserDTO

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch company id from jwt"), nil)
		return
	}

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	if err := c.userService.EditUser(ctx, userId, companyId, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "User updated successfully")

}
