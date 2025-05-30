package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type FetchPayment struct {
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	UserId    string `json:"user_id"`
	CompanyID string `json:"company_id"`
}

type PaymentResponse struct {
	ID             string               `json:"id"`
	CompanyID      string               `json:"company_id"`
	UserID         string               `json:"user_id"`
	SubscriptionID string               `json:"subscription_id"`
	PaymentID      *string              `json:"payment_id"`
	Amount         decimal.Decimal      `json:"amount"`
	Currency       *string              `json:"currency"`
	PaymentMethod  *string              `json:"payment_method"`
	Status         *string              `json:"status"`
	Notes          *string              `json:"notes"`
	CreatedAt      *time.Time           `json:"created_at"`
	UpdatedAt      *time.Time           `json:"updated_at"`
	DeletedAt      *time.Time           `json:"deleted_at"`
	Company        CompanyResponse      `json:"company"`
	Subscription   SubscriptionResponse `json:"subscription"`
	User           UserResponse         `json:"user"`
}

type SubscriptionResponse struct {
	Subscriptionplanid        string          `json:"subscriptionplanid"`
	Subscriptionamount        decimal.Decimal `json:"subscriptionamount"`
	Subscriptionbillingcycle  *string         `json:"subscriptionbillingcycle"`
	Subscriptiontrialstartsat *time.Time      `json:"subscriptiontrialstartsat"`
	Subscriptiontrialendsat   *time.Time      `json:"subscriptiontrialendsat"`
	Subscriptionstartsat      *time.Time      `json:"subscriptionstartsat"`
	Subscriptionendsat        *time.Time      `json:"subscriptionendsat"`
	Subscriptionstatus        *string         `json:"subscriptionstatus"`
	Subscriptioncreatedat     *time.Time      `json:"subscriptioncreatedat"`
}

type CompanyResponse struct {
	Companyname      *string   `json:"companyname"`
	Companycreatedat time.Time `json:"companycreatedat"`
	Companyupdatedat time.Time `json:"companyupdatedat"`
}

type UserResponse struct {
	Userfullname    string     `json:"userfullname"`
	Useremail       string     `json:"useremail"`
	Userphonenumber *string    `json:"userphonenumber"`
	Userpicture     *string    `json:"userpicture"`
	Userverified    bool       `json:"userverified"`
	Userblocked     bool       `json:"userblocked"`
	Userstatus      string     `json:"userstatus"`
	Userlastloginat *time.Time `json:"userlastloginat"`
	Usercreatedat   time.Time  `json:"usercreatedat"`
}
