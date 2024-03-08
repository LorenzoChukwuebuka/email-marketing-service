package model

import (
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
	Status   string
}

type Billing struct {
	Id            int        `json:"-" gorm:"primaryKey"`
	UUID          string     `json:"uuid"`
	UserId        int        `json:"user_id"`
	AmountPaid    float32    `json:"amount_paid" validated:"required"`
	PlanId        int        `json:"plan_id" validated:"required"`
	Email         string     `json:"email"`
	Duration      string     `json:"duration"`
	ExpiryDate    time.Time  `json:"expiry_date"`
	Reference     string     `json:"reference"`
	TransactionId string     `json:"transaction_id"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"type:TIMESTAMP;null;default:null"`
	User          User       `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Plan          Plan       `json:"plan,omitempty" gorm:"foreignKey:PlanId;references:ID"`
}

type BillingResponse struct {
	Id            int          `json:"-"`
	UUID          string       `json:"uuid"`
	UserId        int          `json:"user_id"`
	AmountPaid    float32      `json:"amount_paid"`
	PlanId        int          `json:"plan_id"`
	Duration      string       `json:"duration"`
	ExpiryDate    time.Time    `json:"expiry_date"`
	Reference     string       `json:"reference"`
	TransactionId string       `json:"transaction_id"`
	PaymentMethod string       `json:"payment_method"`
	Status        string       `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
	DeletedAt     string       `json:"deleted_at"`
	User          *UserResponse `json:"user"`
	Plan          *PlanResponse `json:"plan"`
}
