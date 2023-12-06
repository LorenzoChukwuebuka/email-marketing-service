package model

import (
	"database/sql"
	"time"
)

type SubscriptionModel struct {
	Id        int       `json:"id"`
	UUID      string    `json:"uuid"`
	UserId    int       `json:"user_id"`
	PlanId    int       `json:"plan_id"`
	PaymentId int       `json:"payment_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Expired   bool      `json:"expired"`
	TransactionId string `json:"transaction_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	Cancelled bool `json:"cancelled"`
	DateCancelled sql.NullTime `json:"date_cancelled"`
}

type SubscriptionResponseModel struct {
	Id        int       `json:"id"`
	UUID      string    `json:"uuid"`
	UserId    int       `json:"user_id"`
	PlanId    int       `json:"plan_id"`
	PaymentId *int       `json:"payment_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Expired   bool      `json:"expired"`
	TransactionId string `json:"transaction_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Cancelled bool `json:"cancelled"`
	DateCancelled string `json:"date_cancelled"`
	Plan PlanResponse
	User UserResponse
	Billing BillingResponse
}
