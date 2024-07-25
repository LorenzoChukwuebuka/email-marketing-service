package model

import (
	"gorm.io/gorm"
	"time"
)

type Billing struct {
	gorm.Model
	UUID          string    `json:"uuid"`
	UserId        uint      `json:"user_id"`
	AmountPaid    float32   `json:"amount_paid" validated:"required"`
	PlanId        int       `json:"plan_id" validated:"required"`
	Email         string    `json:"email"`
	Duration      string    `json:"duration"`
	ExpiryDate    time.Time `json:"expiry_date"`
	Reference     string    `json:"reference"`
	TransactionId string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Plan          *Plan     `json:"plan,omitempty" gorm:"foreignKey:PlanId;references:ID"`
}

type BillingResponse struct {
	Id            int           `json:"-"`
	UUID          string        `json:"uuid"`
	UserId        uint          `json:"user_id"`
	AmountPaid    float32       `json:"amount_paid"`
	PlanId        int           `json:"plan_id"`
	Duration      string        `json:"duration"`
	ExpiryDate    string        `json:"expiry_date"`
	Reference     string        `json:"reference"`
	TransactionId string        `json:"transaction_id"`
	PaymentMethod string        `json:"payment_method"`
	Status        string        `json:"status"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     *string       `json:"updated_at"`
	DeletedAt     *string       `json:"deleted_at"`
	User          *UserResponse `json:"user,omitempty"`
	Plan          *PlanResponse `json:"plan,omitempty"`
}
