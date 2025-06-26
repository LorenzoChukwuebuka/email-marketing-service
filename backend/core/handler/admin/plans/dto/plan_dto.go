// Updated DTO to better align with your seeder structure
package dto

import (
	"email-marketing-service/internal/enums"
	"github.com/google/uuid"
)

type Plan struct {
	PlanName            string           `json:"plan_name" validate:"required"`
	Description         string           `json:"description" validate:"required"`
	Price               float64          `json:"price"`
	BillingCycle        string           `json:"billing_cycle" validate:"required"` // monthly, yearly, etc.
	Status              enums.PlanStatus `json:"status" validate:"required"`
	Features            []PlanFeature    `json:"features"`
	MailingLimits       MailingLimits    `json:"mailing_limits"`
}

type EditPlan struct {
	UUID                string           `json:"uuid"`
	PlanName            string           `json:"plan_name"`
	Description         string           `json:"description"`
	Price               float64          `json:"price"`
	BillingCycle        string           `json:"billing_cycle"`
	Status              enums.PlanStatus `json:"status"`
	Features            []PlanFeature    `json:"features"`
	MailingLimits       MailingLimits    `json:"mailing_limits"`
}

type PlanFeature struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Value       string `json:"value" validate:"required"` // This stores the feature value/limit
}

type MailingLimits struct {
	DailyLimit           int32 `json:"daily_limit"`
	MonthlyLimit         int32 `json:"monthly_limit"`
	MaxRecipientsPerMail int32 `json:"max_recipients_per_mail"`
}

type PlanResponse struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Price        float64       `json:"price"`
	BillingCycle string        `json:"billing_cycle"`
	Status       string        `json:"status"`
	Features     []PlanFeature `json:"features"`
	MailingLimits MailingLimits `json:"mailing_limits"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
}