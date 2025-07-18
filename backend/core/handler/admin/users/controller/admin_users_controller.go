package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/users/dto"
	"email-marketing-service/core/handler/admin/users/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type AdminUserController struct {
	adminUserServices *services.AdminUsersServices
}

func NewAdminUsersController(adminUserServices *services.AdminUsersServices) *AdminUserController {
	return &AdminUserController{
		adminUserServices: adminUserServices,
	}
}

func (c *AdminUserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	req := &dto.AdminFetchUserDTO{
		Search: search,
		Offset: offset,
		Limit:  limit,
	}

	result, err := c.adminUserServices.GetAllUsers(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminUserController) GetVerifiedUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	req := &dto.AdminFetchUserDTO{
		Search: search,
		Offset: offset,
		Limit:  limit,
	}

	result, err := c.adminUserServices.GetVerifiedUsers(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminUserController) GetUnVerfiedUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	req := &dto.AdminFetchUserDTO{
		Search: search,
		Offset: offset,
		Limit:  limit,
	}

	result, err := c.adminUserServices.GetUnverifiedUsers(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminUserController) BlockUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]

	if err := c.adminUserServices.BlockUser(ctx, userId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "User blocked successfully")
}

func (c *AdminUserController) UnblockUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]

	if err := c.adminUserServices.UnblockUser(ctx, userId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "User unblocked successfully")
}

func (c *AdminUserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]

	if err := c.adminUserServices.VerifyUser(ctx, userId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "User verified successfully")
}

func (c *AdminUserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]

	if err := c.adminUserServices.DeleteUser(ctx, userId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "User deleted successfully")
}

func (c *AdminUserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]

	result, err := c.adminUserServices.GetUserByID(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminUserController) GetUserStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)
	userId := vars["userId"]
	result, err := c.adminUserServices.GetUserStats(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *AdminUserController) SendEmailToUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.AdminEmailLogDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to decode request body"), nil)
		return
	}

	if err := c.adminUserServices.SendEmailToUsers(ctx, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "Email sent successfully")
}

