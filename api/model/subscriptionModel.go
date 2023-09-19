package model

import (
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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriptionResponseModel struct {
	Id        int    `json:"id"`
	UUID      string `json:"uuid"`
	User      UserResponse
	Plan      PlanResponse
	PaymentId PaymentResponse
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Expired   bool      `json:"expired"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
