package adminController

import (
	adminservice "email-marketing-service/api/v1/services/admin"
	"net/http"
	"strconv"
	//"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type AdminUsersController struct {
	AdminUserService *adminservice.AdminUsers
}

// NewAdminUsersController initializes a new AdminUsersController
func NewAdminUsersController(adminUserService *adminservice.AdminUsers) *AdminUsersController {
	return &AdminUsersController{
		AdminUserService: adminUserService,
	}
}

// GetAllUsers retrieves all users with pagination and search
func (c *AdminUsersController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	users, err := c.AdminUserService.GetAllUsers(page, pageSize, searchQuery)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, users)
}

// BlockUser blocks a user by their userId
func (c *AdminUsersController) BlockUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	userResponse, err := c.AdminUserService.BlockUser(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userResponse)
}

// UnblockUser unblocks a user by their userId
func (c *AdminUsersController) UnblockUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	userResponse, err := c.AdminUserService.UnblockUser(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userResponse)
}

// GetVerifiedUsers retrieves all verified users with pagination and search
func (c *AdminUsersController) GetVerifiedUsers(w http.ResponseWriter, r *http.Request) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	users, err := c.AdminUserService.GetVerifiedUsers(page, pageSize, searchQuery)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, users)
}

func (c *AdminUsersController) GetUnVerifiedUsers(w http.ResponseWriter, r *http.Request) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	users, err := c.AdminUserService.GetUnverifiedUsers(page, pageSize, searchQuery)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, users)
}

// GetSingleUser retrieves a single user by their userId
func (c *AdminUsersController) GetSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	userResponse, err := c.AdminUserService.GetSingleUser(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userResponse)
}

// VerifyUser verifies a user by their userId
func (c *AdminUsersController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	userResponse, err := c.AdminUserService.VerifyUser(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userResponse)
}

func (c *AdminUsersController) GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	userResponse, err := c.AdminUserService.GetUserStats(userId)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userResponse)
}
