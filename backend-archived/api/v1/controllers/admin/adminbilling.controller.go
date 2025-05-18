package adminController

import (
	"email-marketing-service/api/v1/services/admin"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AdminBillingController struct {
	AdminBillingSVC *adminservice.AdminBillingService
}

// Constructor for AdminBillingController
func NewAdminBillingController(adminBillingSVC *adminservice.AdminBillingService) *AdminBillingController {
	return &AdminBillingController{
		AdminBillingSVC: adminBillingSVC,
	}
}

// Get total billing (payments) for a user
func (c *AdminBillingController) GetTotalBillingForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]

	// Convert userId string to uint
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		response.ErrorResponse(w, "Invalid user ID")
		return
	}

	// Fetch total billing for the user
	totalBilling, err := c.AdminBillingSVC.GetTotalBillingForUser(uint(userId))
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	// Send the response
	response.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"userId":       userId,
		"totalBilling": totalBilling,
	})
}

// Get all billings (payments) for a user
func (c *AdminBillingController) GetAllBillingsForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]

	// Convert userId string to uint
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		response.ErrorResponse(w, "Invalid user ID")
		return
	}

	// Fetch all billings for the user
	billings, err := c.AdminBillingSVC.GetAllBillingsForUser(uint(userId))
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	// Send the response
	response.SuccessResponse(w, http.StatusOK, billings)
}

// Get total billing for a specific period (1 day, 1 week, 1 month, 1 year)
func (c *AdminBillingController) GetTotalBillingsByPeriod(w http.ResponseWriter, r *http.Request) {
	// Retrieve the period from query parameters
	period := r.URL.Query().Get("period")
	if period == "" {
		response.ErrorResponse(w, "Period is required")
		return
	}

	// Fetch total billings for the specified period
	totalBilling, err := c.AdminBillingSVC.GetTotalBillingsByPeriod(period)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	// Send the response
	response.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"period":       period,
		"totalBilling": totalBilling,
	})
}
