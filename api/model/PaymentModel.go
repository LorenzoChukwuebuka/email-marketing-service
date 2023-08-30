package model

import (
	"database/sql"
	"time"
)

type PaymentModel struct {
	Id         int          `json:"id"`
	UserId     int          `json:"user_id"`
	AmountPaid float32      `json:"amount_paid"`
	PlanId     int          `json:"plan_id"`
	Duration   string       `json:"duration"`
	ExpiryDate time.Time    `json:"expiry_date"`
	Reference  string       `json:"reference"`
	Status     string       `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}
