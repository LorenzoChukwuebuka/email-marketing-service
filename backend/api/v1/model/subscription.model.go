package model

import (
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	gorm.Model
	UUID          string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId        uint       `json:"user_id"`
	PlanId        uint       `json:"plan_id"`
	PaymentId     int        `json:"payment_id"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	Expired       bool       `json:"expired"`
	TransactionId string     `json:"transaction_id"`
	Cancelled     bool       `json:"cancelled"`
	DateCancelled *time.Time `json:"date_cancelled" gorm:"type:TIMESTAMP;null;default:null"`
	Plan          *Plan      `json:"plan"`
	User          *User      `json:"user" gorm:"foreignKey:UserId"`
	Billing       *Billing   `json:"billing,omitempty" gorm:"foreignKey:PaymentId"`
}

type SubscriptionResponseModel struct {
	Id            int              `json:"-"`
	UUID          string           `json:"uuid"`
	UserId        uint             `json:"user_id"`
	PlanId        uint             `json:"plan_id"`
	PaymentId     int              `json:"payment_id"`
	StartDate     time.Time        `json:"start_date"`
	EndDate       time.Time        `json:"end_date"`
	Expired       bool             `json:"expired"`
	TransactionId string           `json:"transaction_id"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
	Cancelled     bool             `json:"cancelled"`
	DateCancelled string           `json:"date_cancelled"`
	Plan          *PlanResponse    `json:"plan"`
	User          *UserResponse    `json:"user"`
	Billing       *BillingResponse `json:"billing"`
}
