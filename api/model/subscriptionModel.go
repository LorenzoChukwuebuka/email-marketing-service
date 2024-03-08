package model

import (
	"time"
)

type Subscription struct {
	Id            int        `json:"-" gorm:"primaryKey"`
	UUID          string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId        int        `json:"user_id"`
	PlanId        int        `json:"plan_id"`
	PaymentId     int        `json:"payment_id"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	Expired       bool       `json:"expired"`
	TransactionId string     `json:"transaction_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:TIMESTAMP;null;default:null"`
	Cancelled     bool       `json:"cancelled"`
	DateCancelled *time.Time `json:"date_cancelled" gorm:"type:TIMESTAMP;null;default:null"`
	Plan          *Plan      `json:"plan,omitempty"`
	User          *User      `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Billing       *Billing   `json:"billing,omitempty" gorm:"foreignKey:PaymentId;references:Id"`
}

type SubscriptionResponseModel struct {
	Id            int       `json:"-"`
	UUID          string    `json:"uuid"`
	UserId        int       `json:"user_id"`
	PlanId        int       `json:"plan_id"`
	PaymentId     *int      `json:"payment_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Expired       bool      `json:"expired"`
	TransactionId string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	Cancelled     bool      `json:"cancelled"`
	DateCancelled string    `json:"date_cancelled"`
	Plan          *PlanResponse
	User          *UserResponse
	Billing       *BillingResponse
}
