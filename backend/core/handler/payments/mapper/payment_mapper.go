package mapper

import (
	"email-marketing-service/core/handler/payments/dto"
	db "email-marketing-service/internal/db/sqlc"
)

func MapPaymentResponse(row db.GetPaymentsByCompanyAndUserRow) *dto.PaymentResponse {
	response := &dto.PaymentResponse{
		ID:             row.ID.String(),
		CompanyID:      row.CompanyID.String(),
		UserID:         row.UserID.String(),
		SubscriptionID: row.SubscriptionID.String(),
		Amount:         row.Amount,
	}

	// Handle nullable PaymentID
	if row.PaymentID.Valid {
		response.PaymentID = &row.PaymentID.String
	}

	// Handle nullable Currency
	if row.Currency.Valid {
		response.Currency = &row.Currency.String
	}

	// Handle nullable PaymentMethod
	if row.PaymentMethod.Valid {
		response.PaymentMethod = &row.PaymentMethod.String
	}

	// Handle nullable Status
	if row.Status.Valid {
		response.Status = &row.Status.String
	}

	// Handle nullable Notes
	if row.Notes.Valid {
		response.Notes = &row.Notes.String
	}

	// Handle nullable CreatedAt
	if row.CreatedAt.Valid {
		response.CreatedAt = &row.CreatedAt.Time
	}

	// Handle nullable UpdatedAt
	if row.UpdatedAt.Valid {
		response.UpdatedAt = &row.UpdatedAt.Time
	}

	// Handle nullable DeletedAt
	if row.DeletedAt.Valid {
		response.DeletedAt = &row.DeletedAt.Time
	}

	// Map Company data
	response.Company = dto.CompanyResponse{
		Companycreatedat: row.Companycreatedat,
		Companyupdatedat: row.Companyupdatedat,
	}

	// Handle nullable Company name
	if row.Companyname.Valid {
		response.Company.Companyname = &row.Companyname.String
	}

	// Map User data
	response.User = dto.UserResponse{
		Userfullname:  row.Userfullname,
		Useremail:     row.Useremail,
		Userverified:  row.Userverified,
		Userblocked:   row.Userblocked,
		Userstatus:    row.Userstatus,
		Usercreatedat: row.Usercreatedat,
	}

	// Handle nullable User phone number
	if row.Userphonenumber.Valid {
		response.User.Userphonenumber = &row.Userphonenumber.String
	}

	// Handle nullable User picture
	if row.Userpicture.Valid {
		response.User.Userpicture = &row.Userpicture.String
	}

	// Handle nullable User last login
	if row.Userlastloginat.Valid {
		response.User.Userlastloginat = &row.Userlastloginat.Time
	}

	// Map Subscription data
	response.Subscription = dto.SubscriptionResponse{
		Subscriptionplanid: row.Subscriptionplanid.String(),
		Subscriptionamount: row.Subscriptionamount,
	}

	// Handle nullable Subscription billing cycle
	if row.Subscriptionbillingcycle.Valid {
		response.Subscription.Subscriptionbillingcycle = &row.Subscriptionbillingcycle.String
	}

	// Handle nullable Subscription trial starts at
	if row.Subscriptiontrialstartsat.Valid {
		response.Subscription.Subscriptiontrialstartsat = &row.Subscriptiontrialstartsat.Time
	}

	// Handle nullable Subscription trial ends at
	if row.Subscriptiontrialendsat.Valid {
		response.Subscription.Subscriptiontrialendsat = &row.Subscriptiontrialendsat.Time
	}

	// Handle nullable Subscription starts at
	if row.Subscriptionstartsat.Valid {
		response.Subscription.Subscriptionstartsat = &row.Subscriptionstartsat.Time
	}

	// Handle nullable Subscription ends at
	if row.Subscriptionendsat.Valid {
		response.Subscription.Subscriptionendsat = &row.Subscriptionendsat.Time
	}

	// Handle nullable Subscription status
	if row.Subscriptionstatus.Valid {
		response.Subscription.Subscriptionstatus = &row.Subscriptionstatus.String
	}

	// Handle nullable Subscription created at
	if row.Subscriptioncreatedat.Valid {
		response.Subscription.Subscriptioncreatedat = &row.Subscriptioncreatedat.Time
	}

	return response
}

// Helper function to map a slice of rows to a slice of responses
func MapPaymentResponses(rows []db.GetPaymentsByCompanyAndUserRow) []*dto.PaymentResponse {
	responses := make([]*dto.PaymentResponse, len(rows))
	for i, row := range rows {
		responses[i] = MapPaymentResponse(row)
	}
	return responses
}
