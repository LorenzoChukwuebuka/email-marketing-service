package model

import (
	"database/sql"
	"time"
)

type BasePaymentModelData struct {
	UserId        int     `json:"user_id" validate:"required"`
	Email         string  `json:"email"`
	AmountToPay   float64 `json:"amount_to_pay" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Duration      string  `json:"duration" validate:"required"`
	PlanId        int     `json:"plan_id" validate:"required"`
}

type BaseProcessPaymentModel struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
	Reference     string `json:"reference" validate:"required"`
}

type BasePaymentResponse struct {
	Amount   float64
	PlanID   int
	UserID   int
	Duration string
	Email    string
	Status string
}

type BillingModel struct {
	Id         int          `json:"id"`
	UUID       string       `json:"uuid"`
	UserId     int          `json:"user_id"`
	AmountPaid float32      `json:"amount_paid" validated:"required"`
	PlanId     int          `json:"plan_id" validated:"required"`
	Email      string       `json:"email"`
	Duration   string       `json:"duration"`
	ExpiryDate time.Time    `json:"expiry_date"`
	Reference  string       `json:"reference"`
	TransactionId string  `json:"transaction_id"`
	PaymentMethod string `json:"payment_method"`
	Status     string       `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}

type PaymentResponse struct {
	Id         int       `json:"id"`
	UUID       string    `json:"uuid"`
	UserId     int       `json:"user_id"`
	AmountPaid float32   `json:"amount_paid"`
	PlanId     int       `json:"plan_id"`
	Duration   string    `json:"duration"`
	ExpiryDate time.Time `json:"expiry_date"`
	Reference  string    `json:"reference"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	DeletedAt  string    `json:"deleted_at"`
	User       UserResponse
	Plan       PlanResponse
}
